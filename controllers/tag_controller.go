// controllers/tag_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllTags godoc
// @Summary Ambil semua tag
// @Description Mengambil semua tag beserta kategorinya
// @Tags Tag
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /tags [get]
func GetAllTags(c *gin.Context) {
	var tags []models.Tag
	
	if err := config.DB.Preload("Category").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// CreateTag godoc
// @Summary Buat tag baru
// @Description Membuat tag baru (admin/editor only)
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.CreateTagInput true "Data tag"
// @Success 201 {object} map[string]interface{}
// @Router /editor/tags [post]
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

	
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "CREATE_TAG", tag.TagID, "tags")

	c.JSON(http.StatusCreated, gin.H{"data": tag})
}

// UpdateTag godoc
// @Summary Update tag
// @Description Mengubah tag yang ada (admin/editor only)
// @Tags Tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID tag"
// @Param body body models.UpdateTagInput true "Data tag"
// @Success 200 {object} map[string]interface{}
// @Router /editor/tags/{id} [put]
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
	
	
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "UPDATE_TAG", tag.TagID, "tags")

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// DeleteTag godoc
// @Summary Hapus tag
// @Description Menghapus tag (admin/editor only)
// @Tags Tag
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID tag"
// @Success 200 {object} map[string]string
// @Router /editor/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	var tag models.Tag
	if err := config.DB.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag tidak ditemukan"})
		return
	}
	
	tagIDtoLog := tag.TagID 
	if err := config.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus tag"})
		return
	}

	
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "DELETE_TAG", tagIDtoLog, "tags")

	c.JSON(http.StatusOK, gin.H{"message": "Tag berhasil dihapus"})
}