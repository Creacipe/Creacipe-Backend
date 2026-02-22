// controllers/menu_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateMenu godoc
// @Summary Buat resep baru
// @Description Membuat resep baru dengan upload gambar (multipart/form-data)
// @Tags Menu
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Judul resep"
// @Param description formData string true "Deskripsi resep"
// @Param ingredients formData string true "Bahan-bahan (JSON)"
// @Param instructions formData string true "Langkah-langkah (JSON)"
// @Param image_file formData file true "Gambar resep"
// @Param tag_ids formData string false "Tag IDs (comma separated, misal: 1,2,3)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /menus [post]
func CreateMenu(c *gin.Context) {
	// 2. Baca form-data
	title := c.PostForm("title")
	description := c.PostForm("description")
	ingredientsJSON := c.PostForm("ingredients")   
	instructionsJSON := c.PostForm("instructions") 

	user := c.MustGet("user").(models.User)
	var finalImageURL string 

	// 3. Proses upload file 
	file, err := c.FormFile("image_file")

	if err != nil {
		// Jika tidak ada file diupload, return error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar wajib diupload"})
		return
	}

	// Upload file ke ImageKit
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca file"})
		return
	}
	defer fileContent.Close()

	imageURL, err := helpers.UploadToImageKit(fileContent, file, "menus")
	if err != nil {
		log.Printf("ImageKit upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke ImageKit"})
		return
	}
	finalImageURL = imageURL

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

	
	helpers.CreateLog(user.UserID, "CREATE_MENU", menu.MenuID, "menus")
	

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// GetAllMenus godoc
// @Summary Ambil semua resep
// @Description Mengambil daftar resep yang sudah disetujui dengan pagination dan search
// @Tags Menu
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param limit query int false "Jumlah per halaman" default(12)
// @Param search query string false "Cari berdasarkan judul"
// @Success 200 {object} map[string]interface{}
// @Router /menus [get]
func GetAllMenus(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build base query
	baseQuery := config.DB.Model(&models.Menu{}).Where("status = ?", "approved")
	
	// Apply search filter if provided (case-insensitive)
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		baseQuery = baseQuery.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
			searchLower, searchLower,
		)
	}

	// Get total count with search filter
	var total int64
	baseQuery.Count(&total)

	// Fetch menus with pagination and search
	var menus []models.Menu
	fetchQuery := config.DB.Preload("User").Preload("Tags").
		Where("status = ?", "approved")
	
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		fetchQuery = fetchQuery.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
			searchLower, searchLower,
		)
	}
	
	if err := fetchQuery.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}

	// Collect menu IDs for batch query
	menuIDs := make([]uint, len(menus))
	for i, menu := range menus {
		menuIDs[i] = menu.MenuID
	}

	// Batch query for votes stats
	type VoteStat struct {
		MenuID        uint
		TotalLikes    int
		TotalDislikes int
	}
	var voteStats []VoteStat
	if len(menuIDs) > 0 {
		config.DB.Table("menu_votes").
			Select("menu_id, IFNULL(SUM(likes_count), 0) as total_likes, IFNULL(SUM(dislikes_count), 0) as total_dislikes").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&voteStats)
	}
	
	// Batch query for bookmark counts
	type BookmarkStat struct {
		MenuID         uint
		TotalBookmarks int
	}
	var bookmarkStats []BookmarkStat
	if len(menuIDs) > 0 {
		config.DB.Table("user_bookmarks").
			Select("menu_id, COUNT(*) as total_bookmarks").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&bookmarkStats)
	}

	// Create maps for quick lookup
	voteMap := make(map[uint]VoteStat)
	for _, vs := range voteStats {
		voteMap[vs.MenuID] = vs
	}
	bookmarkMap := make(map[uint]int)
	for _, bs := range bookmarkStats {
		bookmarkMap[bs.MenuID] = bs.TotalBookmarks
	}

	// Build results with stats
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}

	var results []MenuWithStats
	for _, menu := range menus {
		vs := voteMap[menu.MenuID]
		results = append(results, MenuWithStats{
			Menu:           menu,
			TotalLikes:     vs.TotalLikes,
			TotalDislikes:  vs.TotalDislikes,
			TotalBookmarks: bookmarkMap[menu.MenuID],
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

// GetMenuByID godoc
// @Summary Detail resep
// @Description Mengambil detail satu resep berdasarkan ID
// @Tags Menu
// @Produce json
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [get]
func GetMenuByID(c *gin.Context) {
	var menu models.Menu
	
	// Ambil menu tanpa filter status dulu, dengan preload Tags
	if err := config.DB.Preload("User").Preload("Tags").First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	
	// Hitung total likes, dislikes, dan bookmarks 
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
	
	// Kembalikan data menu 
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateMenu godoc
// @Summary Update resep
// @Description Mengupdate data resep (hanya pemilik atau moderator)
// @Tags Menu
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Param title formData string false "Judul resep"
// @Param description formData string false "Deskripsi"
// @Param ingredients formData string false "Bahan-bahan"
// @Param instructions formData string false "Langkah-langkah"
// @Param image_file formData file false "Gambar baru"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.Preload("Tags").First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	
	
	user := c.MustGet("user").(models.User)


	if user.UserID != menu.UserID && user.Role.RoleName != "admin" && user.Role.RoleName != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah resep ini"})
		return
	}

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
		// Ada file baru diupload - upload ke ImageKit
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca file"})
			return
		}
		defer fileContent.Close()

		imageURL, err := helpers.UploadToImageKit(fileContent, file, "menus")
		if err != nil {
			log.Printf("ImageKit upload error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke ImageKit"})
			return
		}
		menu.ImageURL = imageURL
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

	
	helpers.CreateLog(user.UserID, "UPDATE_MENU", menu.MenuID, "menus")
	
	config.DB.Preload("User").Preload("Tags").First(&menu, menu.MenuID)


	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// DeleteMenu godoc
// @Summary Hapus resep
// @Description Menghapus resep (hanya pemilik atau moderator)
// @Tags Menu
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /menus/{id} [delete]
func DeleteMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	user := c.MustGet("user").(models.User)
	menuIDtoLog := menu.MenuID 

	if user.UserID != menu.UserID && user.Role.RoleName != "admin" && user.Role.RoleName != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus resep ini"})
		return
	}

	if err := config.DB.Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus resep"})
		return
	}

	
	helpers.CreateLog(user.UserID, "DELETE_MENU", menuIDtoLog, "menus")
	

	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil dihapus"})
}

// GetPopularMenus godoc
// @Summary Resep terpopuler
// @Description Mengambil daftar resep terpopuler berdasarkan vote score
// @Tags Menu
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param limit query int false "Jumlah per halaman" default(12)
// @Success 200 {object} map[string]interface{}
// @Router /menus/popular [get]
func GetPopularMenus(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	
	var total int64
	config.DB.Model(&models.Menu{}).Where("status = ?", "approved").Count(&total)


	type MenuStat struct {
		MenuID         uint `gorm:"column:menu_id"`
		TotalLikes     int  `gorm:"column:total_likes"`
		TotalDislikes  int  `gorm:"column:total_dislikes"`
		VoteScore      int  `gorm:"column:vote_score"`
		TotalBookmarks int  `gorm:"column:total_bookmarks"`
	}

	var menuStats []MenuStat
	statsQuery := `
		SELECT 
			m.menu_id,
			COALESCE(SUM(mv.likes_count), 0) as total_likes,
			COALESCE(SUM(mv.dislikes_count), 0) as total_dislikes,
			COALESCE(SUM(mv.likes_count), 0) - COALESCE(SUM(mv.dislikes_count), 0) as vote_score,
			(SELECT COUNT(*) FROM user_bookmarks ub WHERE ub.menu_id = m.menu_id) as total_bookmarks
		FROM menus m
		LEFT JOIN menu_votes mv ON m.menu_id = mv.menu_id
		WHERE m.status = 'approved'
		GROUP BY m.menu_id
		ORDER BY total_likes DESC, vote_score DESC, total_bookmarks DESC
		LIMIT ? OFFSET ?
	`

	if err := config.DB.Raw(statsQuery, limit, offset).Scan(&menuStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep populer"})
		return
	}

	if len(menuStats) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": []interface{}{},
			"meta": gin.H{"total": total, "page": page, "limit": limit, "total_pages": 0},
		})
		return
	}

	
	menuIDs := make([]uint, len(menuStats))
	for i, ms := range menuStats {
		menuIDs[i] = ms.MenuID
	}

	var menus []models.Menu
	if err := config.DB.Preload("Tags").Where("menu_id IN ?", menuIDs).Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail resep"})
		return
	}

	
	menuMap := make(map[uint]models.Menu)
	for _, menu := range menus {
		menuMap[menu.MenuID] = menu
	}

	
	type PopularResult struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
		VoteScore      int `json:"vote_score"`
	}

	var results []PopularResult
	for _, ms := range menuStats {
		if menu, ok := menuMap[ms.MenuID]; ok {
			results = append(results, PopularResult{
				Menu:           menu,
				TotalLikes:     ms.TotalLikes,
				TotalDislikes:  ms.TotalDislikes,
				TotalBookmarks: ms.TotalBookmarks,
				VoteScore:      ms.VoteScore,
			})
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

// GetMyMenus godoc
// @Summary Resep saya
// @Description Mengambil semua resep yang dibuat oleh user yang login
// @Tags Menu
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman" default(1)
// @Param limit query int false "Jumlah per halaman" default(12)
// @Param search query string false "Cari berdasarkan judul"
// @Success 200 {object} map[string]interface{}
// @Router /me/menus [get]
func GetMyMenus(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := strings.TrimSpace(c.Query("search"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 12
	}
	offset := (page - 1) * limit

	
	query := config.DB.Model(&models.Menu{}).Where("user_id = ?", user.UserID)
	

	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchLower, searchLower)
	}
	
	// Get total count
	var total int64
	query.Count(&total)

	var menus []models.Menu
	
	// Build menu query with pagination
	menuQuery := config.DB.Preload("Tags").Where("user_id = ?", user.UserID)
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		menuQuery = menuQuery.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchLower, searchLower)
	}
	
	if err := menuQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}
	
	// Collect menu IDs for batch stats query
	var menuIDs []uint
	for _, menu := range menus {
		menuIDs = append(menuIDs, menu.MenuID)
	}
	
	// Batch query for votes stats
	type VoteStat struct {
		MenuID        uint
		TotalLikes    int
		TotalDislikes int
	}
	var voteStats []VoteStat
	if len(menuIDs) > 0 {
		config.DB.Table("menu_votes").
			Select("menu_id, IFNULL(SUM(likes_count), 0) as total_likes, IFNULL(SUM(dislikes_count), 0) as total_dislikes").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&voteStats)
	}

	// Batch query for bookmark counts
	type BookmarkStat struct {
		MenuID         uint
		TotalBookmarks int
	}
	var bookmarkStats []BookmarkStat
	if len(menuIDs) > 0 {
		config.DB.Table("user_bookmarks").
			Select("menu_id, COUNT(*) as total_bookmarks").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&bookmarkStats)
	}

	// Create maps for quick lookup
	voteMap := make(map[uint]VoteStat)
	for _, v := range voteStats {
		voteMap[v.MenuID] = v
	}
	bookmarkMap := make(map[uint]int)
	for _, b := range bookmarkStats {
		bookmarkMap[b.MenuID] = b.TotalBookmarks
	}
	
	// Build results with stats
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	for _, menu := range menus {
		vs := voteMap[menu.MenuID]
		results = append(results, MenuWithStats{
			Menu:           menu,
			TotalLikes:     vs.TotalLikes,
			TotalDislikes:  vs.TotalDislikes,
			TotalBookmarks: bookmarkMap[menu.MenuID],
		})
	}
	
	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":        results,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetMyBookmarks godoc
// @Summary Bookmark saya
// @Description Mengambil semua resep yang di-bookmark oleh user yang login
// @Tags Menu
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman" default(1)
// @Param limit query int false "Jumlah per halaman" default(12)
// @Param search query string false "Cari berdasarkan judul"
// @Success 200 {object} map[string]interface{}
// @Router /me/bookmarks [get]
func GetMyBookmarks(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := strings.TrimSpace(c.Query("search"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 12
	}
	offset := (page - 1) * limit

	// Get menu IDs yang di-bookmark
	var bookmarkedMenuIDs []uint
	config.DB.Table("user_bookmarks").Where("user_id = ?", user.UserID).Pluck("menu_id", &bookmarkedMenuIDs)
	
	if len(bookmarkedMenuIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":        []interface{}{},
			"total":       0,
			"page":        page,
			"limit":       limit,
			"total_pages": 0,
		})
		return
	}
	
	// Build query for count
	countQuery := config.DB.Model(&models.Menu{}).Where("menu_id IN ?", bookmarkedMenuIDs)
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		countQuery = countQuery.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchLower, searchLower)
	}
	
	var total int64
	countQuery.Count(&total)

	var menus []models.Menu
	
	// Build menu query with pagination
	menuQuery := config.DB.Preload("Tags").Where("menu_id IN ?", bookmarkedMenuIDs)
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		menuQuery = menuQuery.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchLower, searchLower)
	}
	
	if err := menuQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bookmark"})
		return
	}
	
	// Collect menu IDs for batch stats query
	var menuIDs []uint
	for _, menu := range menus {
		menuIDs = append(menuIDs, menu.MenuID)
	}
	
	// Batch query for votes stats
	type VoteStat struct {
		MenuID        uint
		TotalLikes    int
		TotalDislikes int
	}
	var voteStats []VoteStat
	if len(menuIDs) > 0 {
		config.DB.Table("menu_votes").
			Select("menu_id, IFNULL(SUM(likes_count), 0) as total_likes, IFNULL(SUM(dislikes_count), 0) as total_dislikes").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&voteStats)
	}

	// Batch query for bookmark counts
	type BookmarkStat struct {
		MenuID         uint
		TotalBookmarks int
	}
	var bookmarkStats []BookmarkStat
	if len(menuIDs) > 0 {
		config.DB.Table("user_bookmarks").
			Select("menu_id, COUNT(*) as total_bookmarks").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&bookmarkStats)
	}

	// Create maps for quick lookup
	voteMap := make(map[uint]VoteStat)
	for _, v := range voteStats {
		voteMap[v.MenuID] = v
	}
	bookmarkMap := make(map[uint]int)
	for _, b := range bookmarkStats {
		bookmarkMap[b.MenuID] = b.TotalBookmarks
	}
	
	// Build results with stats
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	for _, menu := range menus {
		vs := voteMap[menu.MenuID]
		results = append(results, MenuWithStats{
			Menu:           menu,
			TotalLikes:     vs.TotalLikes,
			TotalDislikes:  vs.TotalDislikes,
			TotalBookmarks: bookmarkMap[menu.MenuID],
		})
	}
	
	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":        results,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetMyCollection godoc
// @Summary Koleksi saya
// @Description Gabungan resep yang dibuat + bookmark user dengan pagination
// @Tags Menu
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman" default(1)
// @Param limit query int false "Jumlah per halaman" default(12)
// @Param search query string false "Cari berdasarkan judul"
// @Success 200 {object} map[string]interface{}
// @Router /me/collection [get]
func GetMyCollection(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	
	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := strings.TrimSpace(c.Query("search"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 12
	}
	offset := (page - 1) * limit
	
	// Get menu IDs dari resep yang dibuat user
	var myMenuIDs []uint
	config.DB.Table("menus").Where("user_id = ?", user.UserID).Pluck("menu_id", &myMenuIDs)
	
	// Get menu IDs dari bookmark
	var bookmarkedMenuIDs []uint
	config.DB.Table("user_bookmarks").Where("user_id = ?", user.UserID).Pluck("menu_id", &bookmarkedMenuIDs)
	
	// Get menu IDs dari resep yang di-like 
	var likedMenuIDs []uint
	config.DB.Table("menu_votes").Where("user_id = ? AND likes_count > 0", user.UserID).Pluck("menu_id", &likedMenuIDs)
	
	// Gabungkan dan deduplikasi IDs 
	menuIDMap := make(map[uint]bool)
	for _, id := range myMenuIDs {
		menuIDMap[id] = true
	}
	for _, id := range bookmarkedMenuIDs {
		menuIDMap[id] = true
	}
	for _, id := range likedMenuIDs {
		menuIDMap[id] = true
	}
	
	var allMenuIDs []uint
	for id := range menuIDMap {
		allMenuIDs = append(allMenuIDs, id)
	}
	
	if len(allMenuIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":        []interface{}{},
			"total":       0,
			"page":        page,
			"limit":       limit,
			"total_pages": 0,
		})
		return
	}
	
	
	countQuery := config.DB.Model(&models.Menu{}).Where("menu_id IN ?", allMenuIDs)
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		countQuery = countQuery.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchLower, searchLower)
	}
	
	var total int64
	countQuery.Count(&total)

	var menus []models.Menu
	
	
	menuQuery := config.DB.Preload("Tags").Where("menu_id IN ?", allMenuIDs)
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		menuQuery = menuQuery.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchLower, searchLower)
	}
	
	if err := menuQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil koleksi"})
		return
	}
	
	
	var menuIDs []uint
	for _, menu := range menus {
		menuIDs = append(menuIDs, menu.MenuID)
	}

	type VoteStat struct {
		MenuID        uint
		TotalLikes    int
		TotalDislikes int
	}
	var voteStats []VoteStat
	if len(menuIDs) > 0 {
		config.DB.Table("menu_votes").
			Select("menu_id, IFNULL(SUM(likes_count), 0) as total_likes, IFNULL(SUM(dislikes_count), 0) as total_dislikes").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&voteStats)
	}

	
	type BookmarkStat struct {
		MenuID         uint
		TotalBookmarks int
	}
	var bookmarkStats []BookmarkStat
	if len(menuIDs) > 0 {
		config.DB.Table("user_bookmarks").
			Select("menu_id, COUNT(*) as total_bookmarks").
			Where("menu_id IN ?", menuIDs).
			Group("menu_id").
			Scan(&bookmarkStats)
	}

	
	voteMap := make(map[uint]VoteStat)
	for _, v := range voteStats {
		voteMap[v.MenuID] = v
	}
	bookmarkMap := make(map[uint]int)
	for _, b := range bookmarkStats {
		bookmarkMap[b.MenuID] = b.TotalBookmarks
	}
	

	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var results []MenuWithStats
	for _, menu := range menus {
		vs := voteMap[menu.MenuID]
		results = append(results, MenuWithStats{
			Menu:           menu,
			TotalLikes:     vs.TotalLikes,
			TotalDislikes:  vs.TotalDislikes,
			TotalBookmarks: bookmarkMap[menu.MenuID],
		})
	}
	
	
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":        results,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}