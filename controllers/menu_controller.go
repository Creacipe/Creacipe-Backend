// controllers/menu_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers" // <-- IMPORT HELPER
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// -----------------------------------------
// CreateMenu menangani pembuatan resep baru oleh member yang login.
func CreateMenu(c *gin.Context) {
	var input models.CreateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.User)

	menu := models.Menu{
		UserID:       user.UserID,
		Title:        input.Title,
		Description:  input.Description,
		Ingredients:  input.Ingredients,
		Instructions: input.Instructions,
		ImageURL:     input.ImageURL,
		Status:       "pending",
	}

	if err := config.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan resep"})
		return
	}

	// Jika ada TagID yang dikirim, hubungkan tag-tag tersebut ke menu
	if len(input.TagIDs) > 0 {
		var tags []models.Tag
		// Cari semua tag berdasarkan ID yang diberikan
		config.DB.Where(input.TagIDs).Find(&tags)
		// Tempelkan tag-tag yang ditemukan ke resep
		config.DB.Model(&menu).Association("Tags").Append(&tags)
	}

	// --- TAMBAHKAN LOG ---
	helpers.CreateLog(user.UserID, "CREATE_MENU", menu.MenuID, "menus")
	// ---------------------

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
	if err := config.DB.Preload("User").Where("status = ?", "approved").First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan atau belum disetujui"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// UpdateMenu menangani pembaruan data resep oleh pemiliknya atau moderator.
func UpdateMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	var input models.UpdateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	user := c.MustGet("user").(models.User)

	// --- PERBAIKAN HAK AKSES ---
	// Pengecekan menggunakan nama peran yang lebih mudah dibaca.
	if user.UserID != menu.UserID && user.Role.RoleName != "admin" && user.Role.RoleName != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah resep ini"})
		return
	}

	if err := config.DB.Model(&menu).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui resep"})
		return
	}

	// --- TAMBAHKAN LOG ---
	helpers.CreateLog(user.UserID, "UPDATE_MENU", menu.MenuID, "menus")
	// ---------------------

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
		TotalBookmarks int    `json:"total_bookmarks"`
		VoteScore      int    `json:"vote_score"`
	}
	var results []PopularResult
	
	// --- PERBAIKAN QUERY ---
	// Query diperbarui untuk menggunakan menu_votes, bukan menu_ratings.
	query := `
        SELECT m.menu_id, m.title, m.image_url, 
        (SELECT COUNT(*) FROM user_bookmarks WHERE menu_id = m.menu_id) as total_bookmarks,
        IFNULL((SELECT SUM(vote_type) FROM menu_votes WHERE menu_id = m.menu_id), 0) as vote_score
        FROM menus m
        WHERE m.status = 'approved'
        ORDER BY vote_score DESC, total_bookmarks DESC
        LIMIT 10
    `
	
	if err := config.DB.Raw(query).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep populer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}