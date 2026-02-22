// controllers/category_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Summary Buat kategori baru
// @Description Membuat kategori tag baru (admin/editor only)
// @Tags Kategori
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.CreateCategoryInput true "Data kategori"
// @Success 201 {object} map[string]interface{}
// @Router /editor/categories [post]
func CreateCategory(c *gin.Context) {
	var input models.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{CategoryName: input.CategoryName}
	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori"})
		return
	}

	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "CREATE_CATEGORY", category.CategoryID, "categories")

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

// GetAllCategories godoc
// @Summary Ambil semua kategori
// @Description Mengambil semua kategori tag
// @Tags Kategori
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /categories [get]
func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// UpdateCategory godoc
// @Summary Update kategori
// @Description Mengubah kategori yang ada (admin/editor only)
// @Tags Kategori
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID kategori"
// @Param body body models.UpdateCategoryInput true "Data kategori"
// @Success 200 {object} map[string]interface{}
// @Router /editor/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	var category models.Category
	if err := config.DB.First(&category, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	var input models.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&category).Update("category_name", input.CategoryName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori"})
		return
	}
	
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "UPDATE_CATEGORY", category.CategoryID, "categories")

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory godoc
// @Summary Hapus kategori
// @Description Menghapus kategori (admin/editor only)
// @Tags Kategori
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID kategori"
// @Success 200 {object} map[string]string
// @Router /editor/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	var category models.Category
	if err := config.DB.First(&category, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	
	categoryIDtoLog := category.CategoryID

	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori, pastikan tidak ada tag yang menggunakannya"})
		return
	}

	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "DELETE_CATEGORY", categoryIDtoLog, "categories")

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}