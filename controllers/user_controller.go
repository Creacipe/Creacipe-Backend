// controllers/user_controller.go
package controllers

import (
	"creacipe-backend/config" // Ganti dengan nama modul Anda
	"creacipe-backend/models" // Ganti dengan nama modul Anda
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsers berfungsi untuk mendapatkan daftar semua pengguna.
// Hanya bisa diakses oleh admin/editor.
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// AdminUpdateUserInput mendefinisikan data user yang bisa diubah oleh admin.
type AdminUpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}

// UpdateUser menangani perubahan data pengguna (nama/email) oleh admin.
func UpdateUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var input AdminUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}



// Struct untuk input saat mengubah peran user
type UpdateUserRoleInput struct {
	Role string `json:"role" binding:"required,oneof=admin editor member"`
}

// UpdateUserRole berfungsi untuk mengubah peran seorang pengguna.
// Hanya bisa diakses oleh admin.
func UpdateUserRole(c *gin.Context) {
	// Ambil user dari database berdasarkan ID di URL
	var user models.User
	userID := c.Param("id")
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Validasi input
	var input UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update peran user
	if err := config.DB.Model(&user).Update("role", input.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui peran pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser menangani penghapusan pengguna oleh admin.
func DeleteUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}