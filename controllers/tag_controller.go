// controllers/tag_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllTags menampilkan semua tag beserta kategorinya.
func GetAllTags(c *gin.Context) {
	var tags []models.Tag
	// Gunakan Preload("Category") untuk menyertakan detail kategori pada setiap tag.
	if err := config.DB.Preload("Category").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// CreateTag membuat tag baru (hanya untuk admin/editor).
func CreateTag(c *gin.Context) {
	var input models.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := models.Tag{
		TagName:    input.TagName,
		CategoryID: input.CategoryID,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat tag"})
		return
	}

	// Catat aktivitas
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "CREATE_TAG", tag.TagID, "tags")

	c.JSON(http.StatusCreated, gin.H{"data": tag})
}

// UpdateTag mengubah tag yang ada (hanya untuk admin/editor).
func UpdateTag(c *gin.Context) {
	var tag models.Tag
	if err := config.DB.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag tidak ditemukan"})
		return
	}

	var input models.UpdateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&tag).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui tag"})
		return
	}
	
	// Catat aktivitas
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "UPDATE_TAG", tag.TagID, "tags")

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// DeleteTag menghapus tag (hanya untuk admin/editor).
func DeleteTag(c *gin.Context) {
	var tag models.Tag
	if err := config.DB.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag tidak ditemukan"})
		return
	}
	
	tagIDtoLog := tag.TagID // Simpan ID untuk log

	if err := config.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus tag"})
		return
	}

	// Catat aktivitas
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "DELETE_TAG", tagIDtoLog, "tags")

	c.JSON(http.StatusOK, gin.H{"message": "Tag berhasil dihapus"})
}