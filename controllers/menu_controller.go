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
	if err := config.DB.Preload("User").Where("status = ?", "approved").Find(&menus).Error; err != nil {
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
	// Kembalikan data menu (frontend yang akan handle logic approved/ownership)
	c.JSON(http.StatusOK, gin.H{"data": menu})
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
	type PopularResult struct {
		MenuID         uint   `json:"menu_id"`
		Title          string `json:"title"`
		ImageURL       string `json:"image_url"`
		TotalLikes     int    `json:"total_likes"`
		TotalDislikes  int    `json:"total_dislikes"`
		TotalBookmarks int    `json:"total_bookmarks"`
		VoteScore      int    `json:"vote_score"`
	}
	var results []PopularResult
	
	// --- PERBAIKAN QUERY untuk struktur baru (likes_count dan dislikes_count) ---
	query := `
        SELECT 
            m.menu_id, 
            m.title, 
            m.image_url, 
            IFNULL(SUM(mv.likes_count), 0) as total_likes,
            IFNULL(SUM(mv.dislikes_count), 0) as total_dislikes,
            (SELECT COUNT(*) FROM user_bookmarks WHERE menu_id = m.menu_id) as total_bookmarks,
            IFNULL(SUM(mv.likes_count) - SUM(mv.dislikes_count), 0) as vote_score
        FROM menus m
        LEFT JOIN menu_votes mv ON m.menu_id = mv.menu_id
        WHERE m.status = 'approved'
        GROUP BY m.menu_id, m.title, m.image_url
        ORDER BY total_likes DESC, vote_score DESC, total_bookmarks DESC
        LIMIT 10
    `
	
	if err := config.DB.Raw(query).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep populer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetMyMenus mendapatkan semua resep yang dibuat oleh user yang sedang login
func GetMyMenus(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	query := `
		SELECT 
			m.*,
			IFNULL(SUM(mv.likes_count), 0) as total_likes,
			IFNULL(SUM(mv.dislikes_count), 0) as total_dislikes,
			(SELECT COUNT(*) FROM user_bookmarks WHERE menu_id = m.menu_id) as total_bookmarks
		FROM menus m
		LEFT JOIN menu_votes mv ON m.menu_id = mv.menu_id
		WHERE m.user_id = ?
		GROUP BY m.menu_id
		ORDER BY m.created_at DESC
	`
	
	if err := config.DB.Raw(query, user.UserID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetMyBookmarks mendapatkan semua resep yang di-bookmark oleh user
func GetMyBookmarks(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	query := `
		SELECT 
			m.menu_id,
			m.user_id,
			m.title,
			m.description,
			m.ingredients,
			m.instructions,
			m.image_url,
			m.status,
			m.created_at,
			m.updated_at,
			COALESCE(SUM(mv.likes_count), 0) as total_likes,
			COALESCE(SUM(mv.dislikes_count), 0) as total_dislikes,
			(SELECT COUNT(*) FROM user_bookmarks WHERE menu_id = m.menu_id) as total_bookmarks
		FROM menus m
		INNER JOIN user_bookmarks ub ON ub.menu_id = m.menu_id
		LEFT JOIN menu_votes mv ON m.menu_id = mv.menu_id
		WHERE ub.user_id = ?
		GROUP BY m.menu_id, m.user_id, m.title, m.description, m.ingredients, m.instructions, 
		         m.image_url, m.status, m.created_at, m.updated_at
		ORDER BY m.created_at DESC
	`
	
	if err := config.DB.Raw(query, user.UserID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bookmark"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetMyCollection mendapatkan gabungan resep yang dibuat user + bookmark
func GetMyCollection(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	query := `
		SELECT DISTINCT 
			m.*,
			IFNULL(SUM(mv.likes_count), 0) as total_likes,
			IFNULL(SUM(mv.dislikes_count), 0) as total_dislikes,
			(SELECT COUNT(*) FROM user_bookmarks WHERE menu_id = m.menu_id) as total_bookmarks
		FROM menus m
		LEFT JOIN menu_votes mv ON m.menu_id = mv.menu_id
		WHERE m.menu_id IN (
			SELECT menu_id FROM menus WHERE user_id = ?
			UNION
			SELECT menu_id FROM user_bookmarks WHERE user_id = ?
		)
		GROUP BY m.menu_id
		ORDER BY m.created_at DESC
	`
	
	if err := config.DB.Raw(query, user.UserID, user.UserID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil koleksi"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": results})
}