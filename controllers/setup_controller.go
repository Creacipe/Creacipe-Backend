// controllers/setup_controller.go
package controllers

import (
	"creacipe-backend/config" // Ganti dengan nama modul Anda
	"creacipe-backend/models" // Ganti dengan nama modul Anda
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SetupInput adalah struct untuk menampung data input setup admin
type SetupInput struct {
	AdminName     string `json:"admin_name" binding:"required"`
	AdminEmail    string `json:"admin_email" binding:"required,email"`
	AdminPassword string `json:"admin_password" binding:"required,min=8"`
}

// SetupFirstAdmin membuat user admin pertama, hanya bisa dijalankan sekali.
func SetupFirstAdmin(c *gin.Context) {
	// 1. Periksa apakah sudah ada admin di database
	var userCount int64
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Setup sudah dijalankan sebelumnya, admin sudah ada."})
		return
	}

	// 2. Validasi input JSON
	var input SetupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// 3. Hash password admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// 4. Buat user admin baru
	admin := models.User{
		Name:     input.AdminName,
		Email:    input.AdminEmail,
		Password: string(hashedPassword),
		Role:     "admin", // Langsung set role sebagai 'admin'
	}

	// 5. Simpan admin ke database
	if result := config.DB.Create(&admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Setup berhasil! Admin pertama telah dibuat."})
}