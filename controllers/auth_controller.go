// controllers/auth_controller.go
package controllers

import (
	"creacipe-backend/config" // Sesuaikan dengan nama modul Anda
	"creacipe-backend/models"
	"creacipe-backend/helpers" // Sesuaikan dengan nama modul Anda
	"net/http"
	"os"
	"time"

	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Fungsi Register menangani logika pendaftaran pengguna baru.
func Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:     3,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// --- 2. TAMBAHKAN LOG UNTUK REGISTRASI ---
	helpers.CreateLog(user.UserID, "USER_REGISTER", user.UserID, "users")
	// ------------------------------------------

	c.JSON(http.StatusOK, gin.H{"message": "Registrasi berhasil"})
}



// Login menangani logika autentikasi pengguna dan pembuatan token JWT.
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau password salah"})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.UserID,
		"role":  user.Role.RoleName,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}
	// --- 3. TAMBAHKAN LOG UNTUK LOGIN ---
	helpers.CreateLog(user.UserID, "USER_LOGIN", user.UserID, "users")
	// ------------------------------------

	c.JSON(http.StatusOK, gin.H{"token": token})
}

//-------untuk reset password-------//
// ForgotPassword membuat token reset.
func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil { /* ... handle error */ }

	var user models.User
	if config.DB.Where("email = ?", body.Email).First(&user).Error != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Jika email terdaftar, instruksi akan dikirim."})
		return
	}

	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	resetToken := hex.EncodeToString(tokenBytes)

	passwordReset := models.PasswordReset{
		UserID:    user.UserID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}
	config.DB.Create(&passwordReset)

	// --- Di sini logika pengiriman email akan ditempatkan ---
	log.Printf("TOKEN RESET UNTUK %s: %s", user.Email, resetToken)
	helpers.CreateLog(user.UserID, "REQUEST_PASSWORD_RESET", user.UserID, "users")
	
	c.JSON(http.StatusOK, gin.H{"message": "Jika email terdaftar, instruksi akan dikirim."})
}

// ResetPassword memvalidasi token dan mengubah password.
func ResetPassword(c *gin.Context) {
	var body struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil { /* ... handle error */ }
	
	var passwordReset models.PasswordReset
	if config.DB.Where("token = ? AND expires_at > ?", body.Token, time.Now()).First(&passwordReset).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token tidak valid atau sudah kedaluwarsa"})
		return
	}

	var user models.User
	config.DB.First(&user, passwordReset.UserID)
	if user.UserID == 0 { /* ... handle error */ }

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	config.DB.Model(&user).Update("password", string(hashedPassword))
	config.DB.Delete(&passwordReset)

	helpers.CreateLog(user.UserID, "COMPLETE_PASSWORD_RESET", user.UserID, "users")
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset."})
}