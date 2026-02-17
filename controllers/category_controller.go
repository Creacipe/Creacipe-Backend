// controllers/category_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCategory membuat kategori tag baru (hanya admin/editor).
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

// GetAllCategories menampilkan semua kategori tag.
func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// UpdateCategory mengubah kategori yang ada (hanya admin/editor).
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

// DeleteCategory menghapus kategori (hanya admin/editor).
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