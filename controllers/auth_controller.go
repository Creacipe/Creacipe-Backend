// controllers/auth_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"creacipe-backend/helpers" 
	"net/http"
	"os"
	"time"
	"fmt"

	// "crypto/rand"
	// "encoding/hex"
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


// Login sekarang mengembalikan access_token dan refresh_token
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mohon lengkapi data login dan captcha"})
		return
	}
	// --- 1. VALIDASI RECAPTCHA (TAMBAHAN BARU) ---
	// Panggil helper yang baru kita buat
	if err := helpers.VerifyRecaptcha(input.RecaptchaToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verifikasi Captcha Gagal: Anda terdeteksi sebagai robot atau token kadaluwarsa."})
		return
	}
	// ---------------------------------------------

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
// RequestPasswordReset mengirim kode verifikasi untuk reset password
func RequestPasswordReset(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password saat ini harus diisi"})
		return
	}

	// Verifikasi password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password saat ini tidak valid"})
		return
	}

	// Generate kode verifikasi 6 digit
	verificationCode := helpers.GenerateVerificationCode()
	expiresAt := time.Now().Add(10 * time.Minute) // Berlaku 10 menit

	// Hapus kode lama yang belum digunakan
	config.DB.Where("user_id = ? AND is_used = false", user.UserID).Delete(&models.PasswordReset{})

	// Simpan kode baru
	passwordReset := models.PasswordReset{
		UserID:           user.UserID,
		VerificationCode: verificationCode,
		ExpiresAt:        expiresAt,
		IsUsed:           false,
	}

	if err := config.DB.Create(&passwordReset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kode verifikasi"})
		return
	}
	
	if err := helpers.SendVerificationEmail(user.Email, user.Name, verificationCode, "reset_password"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal mengirim email verifikasi: %v", err)})
		return
	}
	
	helpers.CreateLog(user.UserID, "REQUEST_PASSWORD_RESET", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Kode verifikasi telah dikirim ke email Anda",
		"expires_at": expiresAt,
	})
}

// VerifyAndResetPassword memverifikasi kode dan reset password
func VerifyAndResetPassword(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var input struct {
		VerificationCode string `json:"verification_code" binding:"required,len=6"`
		NewPassword      string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari kode verifikasi
	var passwordReset models.PasswordReset
	if err := config.DB.Where("user_id = ? AND verification_code = ? AND is_used = false", 
		user.UserID, input.VerificationCode).First(&passwordReset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode verifikasi tidak valid"})
		return
	}

	// Cek apakah kode sudah expired
	if time.Now().After(passwordReset.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode verifikasi sudah kadaluarsa"})
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// Update password
	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah password"})
		return
	}

	// Tandai kode sebagai sudah digunakan
	config.DB.Model(&passwordReset).Update("is_used", true)

	helpers.CreateLog(user.UserID, "PASSWORD_RESET_SUCCESS", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

// ========== EMAIL CHANGE WITH VERIFICATION ==========

// RequestEmailChange mengirim kode verifikasi untuk ubah email
func RequestEmailChange(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var input struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah email sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("email = ? AND user_id != ?", input.NewEmail, user.UserID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan oleh pengguna lain"})
		return
	}

	// Generate kode verifikasi
	verificationCode := helpers.GenerateVerificationCode()
	expiresAt := time.Now().Add(10 * time.Minute)

	// Hapus kode lama yang belum digunakan
	config.DB.Where("user_id = ? AND is_used = false", user.UserID).Delete(&models.EmailVerification{})

	// Simpan kode baru
	emailVerification := models.EmailVerification{
		UserID:           user.UserID,
		NewEmail:         input.NewEmail,
		VerificationCode: verificationCode,
		ExpiresAt:        expiresAt,
		IsUsed:           false,
	}

	if err := config.DB.Create(&emailVerification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kode verifikasi"})
		return
	}

	
	if err := helpers.SendVerificationEmail(input.NewEmail, user.Name, verificationCode, "change_email"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal mengirim email verifikasi: %v", err)})
		return
	}
	
	helpers.CreateLog(user.UserID, "REQUEST_EMAIL_CHANGE", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Kode verifikasi telah dikirim ke email baru Anda",
		"new_email":  input.NewEmail,
		"expires_at": expiresAt,
	})
}

// VerifyAndChangeEmail memverifikasi kode dan ubah email
func VerifyAndChangeEmail(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var input struct {
		VerificationCode string `json:"verification_code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari kode verifikasi
	var emailVerification models.EmailVerification
	if err := config.DB.Where("user_id = ? AND verification_code = ? AND is_used = false", 
		user.UserID, input.VerificationCode).First(&emailVerification).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode verifikasi tidak valid"})
		return
	}

	// Cek apakah kode sudah expired
	if time.Now().After(emailVerification.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode verifikasi sudah kadaluarsa"})
		return
	}

	// Cek lagi apakah email masih available
	var existingUser models.User
	if err := config.DB.Where("email = ? AND user_id != ?", emailVerification.NewEmail, user.UserID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan oleh pengguna lain"})
		return
	}

	// Update email
	if err := config.DB.Model(&user).Update("email", emailVerification.NewEmail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah email"})
		return
	}

	// Tandai kode sebagai sudah digunakan
	config.DB.Model(&emailVerification).Update("is_used", true)

	helpers.CreateLog(user.UserID, "EMAIL_CHANGE_SUCCESS", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Email berhasil diubah",
		"new_email": emailVerification.NewEmail,
	})
}

// ========== FORGOT PASSWORD (PUBLIC - NO AUTH REQUIRED) ==========

// ForgotPasswordRequest mengirim kode OTP ke email user (tanpa perlu login)
func ForgotPasswordRequest(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email harus diisi dengan format yang valid"})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// JANGAN beritahu user bahwa email tidak ditemukan (security)
		// Kirim respons sukses agar attacker tidak tahu email mana yang terdaftar
		c.JSON(http.StatusOK, gin.H{
			"message": "Jika email terdaftar, kode verifikasi akan dikirim ke email Anda",
		})
		return
	}

	// Generate kode verifikasi 6 digit
	verificationCode := helpers.GenerateVerificationCode()
	expiresAt := time.Now().Add(10 * time.Minute) // Berlaku 10 menit

	// Hapus kode lama yang belum digunakan untuk user ini
	config.DB.Where("user_id = ? AND is_used = false", user.UserID).Delete(&models.PasswordReset{})

	// Simpan kode baru
	passwordReset := models.PasswordReset{
		UserID:           user.UserID,
		VerificationCode: verificationCode,
		ExpiresAt:        expiresAt,
		IsUsed:           false,
	}

	if err := config.DB.Create(&passwordReset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kode verifikasi"})
		return
	}

	// Kirim email
	fmt.Printf("[INFO] [FORGOT PASSWORD] Mengirim kode verifikasi ke email: %s\n", user.Email)
	fmt.Printf("[INFO] Kode verifikasi: %s\n", verificationCode)
	
	if err := helpers.SendVerificationEmail(user.Email, user.Name, verificationCode, "reset_password"); err != nil {
		fmt.Printf("[ERROR] Gagal mengirim email: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim email verifikasi"})
		return
	}
	
	fmt.Printf("[SUCCESS] Email lupa password berhasil dikirim ke: %s\n", user.Email)

	helpers.CreateLog(user.UserID, "FORGOT_PASSWORD_REQUEST", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Kode verifikasi telah dikirim ke email Anda. Berlaku 10 menit.",
		"expires_at": expiresAt,
	})
}

// ForgotPasswordVerify memverifikasi kode OTP dan reset password (tanpa perlu login)
func ForgotPasswordVerify(c *gin.Context) {
	var input struct {
		Email            string `json:"email" binding:"required,email"`
		VerificationCode string `json:"verification_code" binding:"required,len=6"`
		NewPassword      string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid. Pastikan semua field diisi dengan benar."})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak terdaftar"})
		return
	}

	// Cari kode verifikasi yang valid
	var passwordReset models.PasswordReset
	if err := config.DB.Where("user_id = ? AND verification_code = ? AND is_used = false", 
		user.UserID, input.VerificationCode).First(&passwordReset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode verifikasi tidak valid"})
		return
	}

	// Cek apakah kode sudah expired
	if time.Now().After(passwordReset.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode verifikasi sudah kadaluarsa. Silakan minta kode baru."})
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// Update password
	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah password"})
		return
	}

	// Tandai kode sebagai sudah digunakan
	config.DB.Model(&passwordReset).Update("is_used", true)

	helpers.CreateLog(user.UserID, "FORGOT_PASSWORD_SUCCESS", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset. Silakan login dengan password baru Anda."})
}