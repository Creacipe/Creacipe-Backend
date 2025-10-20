package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddTagInput mendefinisikan input untuk menambahkan tag ke resep.
type AddTagInput struct {
	TagID uint `json:"tag_id" binding:"required"`
}

// AddTagToMenu menangani logika penambahan tag ke sebuah resep.
func AddTagToMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	var input AddTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}
	var tag models.Tag
	if err := config.DB.First(&tag, input.TagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag tidak ditemukan"})
		return
	}
	// GORM akan menangani penambahan relasi di tabel pivot 'menu_tags'.
	config.DB.Model(&menu).Association("Tags").Append(&tag)
	c.JSON(http.StatusOK, gin.H{"message": "Tag berhasil ditambahkan ke resep"})
}

// BookmarkMenu menangani penambahan resep ke daftar bookmark pengguna.
func BookmarkMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	user := c.MustGet("user").(models.User)
	// GORM akan menangani penambahan relasi di tabel pivot 'user_bookmarks'.
	config.DB.Model(&user).Association("Bookmarks").Append(&menu)
	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil di-bookmark"})
}

// UnbookmarkMenu menangani penghapusan resep dari daftar bookmark.
func UnbookmarkMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	user := c.MustGet("user").(models.User)
	config.DB.Model(&user).Association("Bookmarks").Delete(&menu)
	c.JSON(http.StatusOK, gin.H{"message": "Bookmark berhasil dihapus"})
}

// RateMenuInput mendefinisikan input untuk memberi rating.
type RateMenuInput struct {
	Rating uint `json:"rating" binding:"required,min=1,max=5"`
}

// RateMenu menangani logika pemberian atau pembaruan rating resep.
func RateMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	var input RateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid, rating harus 1-5"})
		return
	}
	user := c.MustGet("user").(models.User)
	
	var existingRating models.MenuRating
	// Cari apakah user sudah pernah memberi rating pada resep ini.
	err := config.DB.Where("user_id = ? AND menu_id = ?", user.UserID, menu.MenuID).First(&existingRating).Error

	if err == nil { // Jika sudah ada, perbarui rating yang ada.
		existingRating.Rating = input.Rating
		config.DB.Save(&existingRating)
	} else { // Jika belum ada, buat entri rating baru.
		newRating := models.MenuRating{UserID: user.UserID, MenuID: menu.MenuID, Rating: input.Rating}
		config.DB.Create(&newRating)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rating berhasil disimpan"})
}