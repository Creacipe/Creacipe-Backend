package tests

import (
	"creacipe-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

// createAuthContext membuat context dengan user yang sudah authenticated
func createAuthContext(user models.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user", user)
		c.Next()
	}
}

// createTestUser membuat user lengkap dengan role dan profile untuk testing
func createTestUser(db *gorm.DB, email string, roleID uint) models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      email,
		Password:   string(hashedPassword),
		RoleID:     roleID,
		StatusUser: "active",
	}
	db.Create(&user)
	db.Where("email = ?", email).First(&user)

	// Create profile for the user
	profile := models.UserProfile{
		UserID: user.UserID,
		Bio:    "Test bio",
	}
	db.Create(&profile)

	// Reload user with relations
	db.Preload("Role").Preload("Profile").Where("email = ?", email).First(&user)
	
	return user
}
