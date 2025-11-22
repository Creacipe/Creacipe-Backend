// controllers/menu_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers" // <-- IMPORT HELPER
	"creacipe-backend/models"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// -----------------------------------------
// -----------------------------------------
// CreateMenu menangani pembuatan resep baru (DENGAN UPLOAD FILE + TAGS)
func CreateMenu(c *gin.Context) {
	// 2. Baca form-data
	title := c.PostForm("title")
	description := c.PostForm("description")
	ingredientsJSON := c.PostForm("ingredients")   // Frontend harus kirim JSON string
	instructionsJSON := c.PostForm("instructions") // Frontend harus kirim JSON string

	user := c.MustGet("user").(models.User)
	var finalImageURL string // Variabel untuk menyimpan URL gambar

	// 3. Proses upload file (WAJIB)
	file, err := c.FormFile("image_file")

	if err != nil {
		// Jika tidak ada file diupload, return error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar wajib diupload"})
		return
	}

	// Upload file ke folder assets
	slug := helpers.Slugify(title)
	ext := filepath.Ext(file.Filename)
	uniqueFilename := helpers.FindUniqueFilename(slug, ext)
	savePath := filepath.Join("assets", uniqueFilename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}
	finalImageURL = "http://localhost:8080/assets/" + uniqueFilename

	// 4. Buat struct Menu untuk disimpan
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        title,
		Description:  description,
		Ingredients:  ingredientsJSON,
		Instructions: instructionsJSON,
		ImageURL:     finalImageURL,
		Status:       "pending",
	}

	if err := config.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan resep ke database"})
		return
	}

	// --- 5. LOGIKA TAGIDS YANG DIPERBAIKI (DITAMBAHKAN KEMBALI) ---
	// Ambil TagIDs sebagai string (misal: "1,3,5")
	tagIDsStr := c.PostForm("tag_ids")
	if tagIDsStr != "" {
		// Ubah string "1,3,5" menjadi array string ["1", "3", "5"]
		tagIDs := strings.Split(tagIDsStr, ",")

		if len(tagIDs) > 0 {
			var tags []models.Tag
			// Cari semua tag berdasarkan array ID
			config.DB.Where(tagIDs).Find(&tags)
			
			// Tempelkan tag-tag yang ditemukan ke resep
			if len(tags) > 0 {
				config.DB.Model(&menu).Association("Tags").Append(&tags)
			}
		}
	}
	// --------------------------------------------------------

	// --- 6. HELPER LOG TETAP DI SINI ---
	helpers.CreateLog(user.UserID, "CREATE_MENU", menu.MenuID, "menus")
	// ------------------------------------

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// GetAllMenus mengambil daftar semua resep yang sudah disetujui.
func GetAllMenus(c *gin.Context) {
	var menus []models.Menu
	if err := config.DB.Preload("User").Preload("Tags").Where("status = ?", "approved").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// GetMenuByID mengambil detail satu resep berdasarkan ID (jika sudah disetujui).
func GetMenuByID(c *gin.Context) {
	var menu models.Menu
	
	// Ambil menu tanpa filter status dulu, dengan preload Tags
	if err := config.DB.Preload("User").Preload("Tags").First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	
	// Hitung total likes, dislikes, dan bookmarks (sama seperti GetPopularMenus)
	var totalLikes, totalDislikes int64
	config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(likes_count), 0)").Scan(&totalLikes)
	config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(dislikes_count), 0)").Scan(&totalDislikes)
	
	var totalBookmarks int64
	config.DB.Table("user_bookmarks").Where("menu_id = ?", menu.MenuID).Count(&totalBookmarks)
	
	// Buat response dengan count
	type MenuDetailResult struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	result := MenuDetailResult{
		Menu:           menu,
		TotalLikes:     int(totalLikes),
		TotalDislikes:  int(totalDislikes),
		TotalBookmarks: int(totalBookmarks),
	}
	
	// Kembalikan data menu (frontend yang akan handle logic approved/ownership)
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateMenu menangani pembaruan data resep oleh pemiliknya atau moderator.
func UpdateMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.Preload("Tags").First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	// var input models.UpdateMenuInput
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	
	user := c.MustGet("user").(models.User)

	// --- PERBAIKAN HAK AKSES ---
	// Pengecekan menggunakan nama peran yang lebih mudah dibaca.
	if user.UserID != menu.UserID && user.Role.RoleName != "admin" && user.Role.RoleName != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah resep ini"})
		return
	}

	// Baca form-data (seperti CreateMenu)
	title := c.PostForm("title")
	description := c.PostForm("description")
	ingredientsJSON := c.PostForm("ingredients")
	instructionsJSON := c.PostForm("instructions")

	// Update field-field text
	if title != "" {
		menu.Title = title
	}
	if description != "" {
		menu.Description = description
	}
	if ingredientsJSON != "" {
		menu.Ingredients = ingredientsJSON
	}
	if instructionsJSON != "" {
		menu.Instructions = instructionsJSON
	}

	// Proses upload file jika ada
	file, err := c.FormFile("image_file")
	if err == nil {
		// Ada file baru diupload
		slug := helpers.Slugify(title)
		ext := filepath.Ext(file.Filename)
		uniqueFilename := helpers.FindUniqueFilename(slug, ext)
		savePath := filepath.Join("assets", uniqueFilename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}
		menu.ImageURL = "http://localhost:8080/assets/" + uniqueFilename
	}

	// Update tags jika ada
	tagIDsStr := c.PostForm("tag_ids")
	if tagIDsStr != "" {
		// Hapus tag lama
		config.DB.Model(&menu).Association("Tags").Clear()

		// Tambah tag baru
		tagIDs := strings.Split(tagIDsStr, ",")
		if len(tagIDs) > 0 {
			var tags []models.Tag
			config.DB.Where(tagIDs).Find(&tags)
			if len(tags) > 0 {
				config.DB.Model(&menu).Association("Tags").Append(&tags)
			}
		}
	}

	// Simpan perubahan
	if err := config.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui resep"})
		return
	}

	// --- TAMBAHKAN LOG ---
	helpers.CreateLog(user.UserID, "UPDATE_MENU", menu.MenuID, "menus")
	// ---------------------
	// Reload data menu dengan relasi
	config.DB.Preload("User").Preload("Tags").First(&menu, menu.MenuID)


	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// DeleteMenu menangani penghapusan resep oleh pemiliknya atau moderator.
func DeleteMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	user := c.MustGet("user").(models.User)
	menuIDtoLog := menu.MenuID // Simpan ID untuk log sebelum dihapus

	// --- PERBAIKAN HAK AKSES ---
	if user.UserID != menu.UserID && user.Role.RoleName != "admin" && user.Role.RoleName != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus resep ini"})
		return
	}

	if err := config.DB.Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus resep"})
		return
	}

	// --- TAMBAHKAN LOG ---
	helpers.CreateLog(user.UserID, "DELETE_MENU", menuIDtoLog, "menus")
	// ---------------------

	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil dihapus"})
}

// GetPopularMenus mengambil daftar resep terpopuler.
func GetPopularMenus(c *gin.Context) {
	// Fetch menus dengan Preload Tags
	var menus []models.Menu
	if err := config.DB.Preload("Tags").Where("status = ?", "approved").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep populer"})
		return
	}

	type PopularResult struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
		VoteScore      int `json:"vote_score"`
	}

	var results []PopularResult
	for _, menu := range menus {
		var totalLikes, totalDislikes int64
		var totalBookmarks int64

		// Get likes and dislikes
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(likes_count), 0)").Scan(&totalLikes)
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(dislikes_count), 0)").Scan(&totalDislikes)

		// Get bookmarks count
		config.DB.Table("user_bookmarks").Where("menu_id = ?", menu.MenuID).Count(&totalBookmarks)

		voteScore := int(totalLikes) - int(totalDislikes)

		results = append(results, PopularResult{
			Menu:           menu,
			TotalLikes:     int(totalLikes),
			TotalDislikes:  int(totalDislikes),
			TotalBookmarks: int(totalBookmarks),
			VoteScore:      voteScore,
		})
	}

	// Sort by total_likes DESC, vote_score DESC, total_bookmarks DESC
	// Simple bubble sort for top 10
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].TotalLikes < results[j+1].TotalLikes ||
				(results[j].TotalLikes == results[j+1].TotalLikes && results[j].VoteScore < results[j+1].VoteScore) ||
				(results[j].TotalLikes == results[j+1].TotalLikes && results[j].VoteScore == results[j+1].VoteScore && results[j].TotalBookmarks < results[j+1].TotalBookmarks) {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	// Limit to top 10
	if len(results) > 10 {
		results = results[:10]
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetMyMenus mendapatkan semua resep yang dibuat oleh user yang sedang login
func GetMyMenus(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var menus []models.Menu
	
	// Gunakan GORM dengan Preload untuk mendapatkan Tags
	if err := config.DB.Preload("Tags").Where("user_id = ?", user.UserID).Order("created_at DESC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}
	
	// Tambahkan stats untuk setiap menu
	
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	for _, menu := range menus {
		var totalLikes, totalDislikes int64
		var totalBookmarks int64
		
		// Get likes and dislikes
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(likes_count), 0)").Scan(&totalLikes)
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(dislikes_count), 0)").Scan(&totalDislikes)
		
		// Get bookmarks count
		config.DB.Table("user_bookmarks").Where("menu_id = ?", menu.MenuID).Count(&totalBookmarks)
		
		results = append(results, MenuWithStats{
			Menu:           menu,
			TotalLikes:     int(totalLikes),
			TotalDislikes:  int(totalDislikes),
			TotalBookmarks: int(totalBookmarks),
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetMyBookmarks mendapatkan semua resep yang di-bookmark oleh user
func GetMyBookmarks(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	// Get menu IDs yang di-bookmark
	var bookmarkedMenuIDs []uint
	config.DB.Table("user_bookmarks").Where("user_id = ?", user.UserID).Pluck("menu_id", &bookmarkedMenuIDs)
	
	if len(bookmarkedMenuIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}
	
	var menus []models.Menu
	
	// Gunakan GORM dengan Preload untuk mendapatkan Tags
	if err := config.DB.Preload("Tags").Where("menu_id IN ?", bookmarkedMenuIDs).Order("created_at DESC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bookmark"})
		return
	}
	
	// Tambahkan stats untuk setiap menu
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	for _, menu := range menus {
		var totalLikes, totalDislikes int64
		var totalBookmarks int64
		
		// Get likes and dislikes
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(likes_count), 0)").Scan(&totalLikes)
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(dislikes_count), 0)").Scan(&totalDislikes)
		
		// Get bookmarks count
		config.DB.Table("user_bookmarks").Where("menu_id = ?", menu.MenuID).Count(&totalBookmarks)
		
		results = append(results, MenuWithStats{
			Menu:           menu,
			TotalLikes:     int(totalLikes),
			TotalDislikes:  int(totalDislikes),
			TotalBookmarks: int(totalBookmarks),
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetMyCollection mendapatkan gabungan resep yang dibuat user + bookmark
func GetMyCollection(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	
	// Get menu IDs dari resep yang dibuat user
	var myMenuIDs []uint
	config.DB.Table("menus").Where("user_id = ?", user.UserID).Pluck("menu_id", &myMenuIDs)
	
	// Get menu IDs dari bookmark
	var bookmarkedMenuIDs []uint
	config.DB.Table("user_bookmarks").Where("user_id = ?", user.UserID).Pluck("menu_id", &bookmarkedMenuIDs)
	
	// Gabungkan dan deduplikasi IDs
	menuIDMap := make(map[uint]bool)
	for _, id := range myMenuIDs {
		menuIDMap[id] = true
	}
	for _, id := range bookmarkedMenuIDs {
		menuIDMap[id] = true
	}
	
	var allMenuIDs []uint
	for id := range menuIDMap {
		allMenuIDs = append(allMenuIDs, id)
	}
	
	if len(allMenuIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}
	
	var menus []models.Menu
	
	// Gunakan GORM dengan Preload untuk mendapatkan Tags
	if err := config.DB.Preload("Tags").Where("menu_id IN ?", allMenuIDs).Order("created_at DESC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil koleksi"})
		return
	}
	
	// Tambahkan stats untuk setiap menu
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	for _, menu := range menus {
		var totalLikes, totalDislikes int64
		var totalBookmarks int64
		
		// Get likes and dislikes
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(likes_count), 0)").Scan(&totalLikes)
		config.DB.Model(&models.MenuVote{}).Where("menu_id = ?", menu.MenuID).Select("IFNULL(SUM(dislikes_count), 0)").Scan(&totalDislikes)
		
		// Get bookmarks count
		config.DB.Table("user_bookmarks").Where("menu_id = ?", menu.MenuID).Count(&totalBookmarks)
		
		results = append(results, MenuWithStats{
			Menu:           menu,
			TotalLikes:     int(totalLikes),
			TotalDislikes:  int(totalDislikes),
			TotalBookmarks: int(totalBookmarks),
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"data": results})
}