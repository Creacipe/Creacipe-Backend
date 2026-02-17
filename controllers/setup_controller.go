// controllers/setup_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers" 
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SetupFirstAdmin membuat user admin pertama, hanya bisa dijalankan sekali.
func SetupFirstAdmin(c *gin.Context) {
	// 1. Periksa apakah sudah ada admin di database menggunakan JOIN.
	var userCount int64
	config.DB.Model(&models.User{}).
		Joins("join roles on roles.role_id = users.role_id").
		Where("roles.role_name = ?", "admin").
		Count(&userCount)
	
	if userCount > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Setup sudah dijalankan sebelumnya, admin sudah ada."})
		return
	}

	// 2. Validasi input JSON.
	var input models.SetupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// 3. Hash password admin.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// 4. Buat user admin baru (mengasumsikan role_id 'admin' adalah 1).
	admin := models.User{
		Name:     input.AdminName,
		Email:    input.AdminEmail,
		Password: string(hashedPassword),
		RoleID:   1, 
	}

	// 5. Simpan admin ke database.
	if result := config.DB.Create(&admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user admin"})
		return
	}

	
	profile := models.UserProfile{UserID: admin.UserID}
	config.DB.Create(&profile)
	

	helpers.CreateLog(admin.UserID, "SETUP_FIRST_ADMIN", admin.UserID, "users")


	c.JSON(http.StatusCreated, gin.H{"message": "Setup berhasil! Admin pertama telah dibuat."})
}