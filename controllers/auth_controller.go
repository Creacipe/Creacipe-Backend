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

	// --- TAMBAHAN: BUAT PROFIL KOSONG OTOMATIS ---
	// Setelah user berhasil dibuat, langsung buatkan profil kosong untuknya.
	profile := models.UserProfile{UserID: user.UserID}
	if err := config.DB.Create(&profile).Error; err != nil {
		// Idealnya, ada penanganan jika pembuatan profil gagal
		log.Printf("Gagal membuat profil untuk user ID %d: %v", user.UserID, err)
	}
	// ---------------------------------------------

	// --- 2. TAMBAHKAN LOG UNTUK REGISTRASI ---
	helpers.CreateLog(user.UserID, "USER_REGISTER", user.UserID, "users")
	// ------------------------------------------

	c.JSON(http.StatusOK, gin.H{"message": "Registrasi berhasil"})
}



// Login menangani logika autentikasi pengguna dan pembuatan token JWT.
// Login menangani logika autentikasi dengan pengecekan status.
// Login sekarang mengembalikan access_token dan refresh_token
// --- FUNGSI LOGIN (DIPERBARUI) ---
// Login sekarang mengembalikan access_token dan refresh_token
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	var user models.User
    // Pastikan Preload("Role") ada di sini
	if err := config.DB.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau password salah"})
		return
	}

	if user.StatusUser == "inactive" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah dinonaktifkan."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau password salah"})
		return
	}

	// 1. Buat Access Token (15 menit)
    // KIRIM datanya secara EKSPLISIT
	accessToken, err := helpers.CreateAccessToken(user.UserID, user.Role.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat access token"})
		return
	}

	// 2. Buat Refresh Token (7 hari)
    // KIRIM datanya secara EKSPLISIT
	refreshToken, err := helpers.CreateRefreshToken(user.UserID, user.Role.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat refresh token"})
		return
	}

	helpers.CreateLog(user.UserID, "USER_LOGIN", user.UserID, "users")
	
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

// --- FUNGSI REFRESH TOKEN (BARU & DIPERBARUI) ---
// RefreshToken memvalidasi refresh token dan memberikan access token baru
func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token dibutuhkan"})
		return
	}

	// 1. Validasi token
	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token tidak valid atau kedaluwarsa"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca token"})
		return
	}

	// 2. Ambil UserID dari token
	userID := uint(claims["sub"].(float64))
	var user models.User
	
	// 3. Ambil data user terbaru (DENGAN ROLE)
	if err := config.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if user.StatusUser == "inactive" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah dinonaktifkan."})
		return
	}

	// 4. Buat Access Token baru
    // KIRIM datanya secara EKSPLISIT
	newAccessToken, err := helpers.CreateAccessToken(user.UserID, user.Role.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat access token baru"})
		return
	}

	// 5. Kembalikan token baru
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
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