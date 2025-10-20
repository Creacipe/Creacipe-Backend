package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateMenuInput mendefinisikan data yang dibutuhkan untuk membuat resep.
type CreateMenuInput struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Ingredients  string `json:"ingredients" binding:"required"`
	Instructions string `json:"instructions" binding:"required"`
	ImageURL     string `json:"image_url"`
}

// CreateMenu menangani pembuatan resep baru oleh member yang login.
func CreateMenu(c *gin.Context) {
	var input CreateMenuInput
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

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// GetAllMenus mengambil daftar semua resep yang sudah disetujui.
func GetAllMenus(c *gin.Context) {
	var menus []models.Menu
	config.DB.Preload("User").Where("status = ?", "approved").Find(&menus)
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

// UpdateMenuInput mendefinisikan data yang bisa diubah pada resep.
type UpdateMenuInput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
	ImageURL     string `json:"image_url"`
}

// UpdateMenu menangani pembaruan data resep oleh pemiliknya atau moderator.
func UpdateMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	var input UpdateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	user := c.MustGet("user").(models.User)
	if user.UserID != menu.UserID && user.Role != "admin" && user.Role != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk mengubah resep ini"})
		return
	}

	config.DB.Model(&menu).Updates(input)
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
	if user.UserID != menu.UserID && user.Role != "admin" && user.Role != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus resep ini"})
		return
	}

	config.DB.Delete(&menu)
	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil dihapus"})
}

// GetPopularMenus mengambil daftar resep terpopuler berdasarkan skor.
func GetPopularMenus(c *gin.Context) {
	// Struct untuk menampung hasil query SQL kustom kita.
	type PopularResult struct {
		MenuID        uint    `json:"menu_id"`
		Title         string  `json:"title"`
		ImageURL      string  `json:"image_url"`
		TotalBookmarks int    `json:"total_bookmarks"`
		AverageRating float64 `json:"average_rating"`
	}

	var results []PopularResult
	
    // Query SQL mentah untuk menghitung skor dan mengurutkan resep.
	query := `
        SELECT m.menu_id, m.title, m.image_url, 
        (SELECT COUNT(*) FROM user_bookmarks WHERE menu_id = m.menu_id) as total_bookmarks,
        IFNULL((SELECT AVG(rating) FROM menu_ratings WHERE menu_id = m.menu_id), 0) as average_rating
        FROM menus m
        WHERE m.status = 'approved'
        ORDER BY total_bookmarks DESC, average_rating DESC
        LIMIT 10
    `
	
	if err := config.DB.Raw(query).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep populer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}