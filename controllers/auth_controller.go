// controllers/auth_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"fmt"
	"net/http"
	"os"
	"time"

	// "crypto/rand"
	// "encoding/hex"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Registrasi user baru
// @Description Mendaftarkan pengguna baru ke sistem
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.RegisterInput true "Data registrasi"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mohon lengkapi data pendaftaran dan captcha"})
		return
	}
	if err := helpers.VerifyRecaptcha(input.RecaptchaToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verifikasi Captcha Gagal: Anda terdeteksi sebagai robot atau token kadaluwarsa."})
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


	profile := models.UserProfile{UserID: user.UserID}
	if err := config.DB.Create(&profile).Error; err != nil {
		log.Printf("Gagal membuat profil untuk user ID %d: %v", user.UserID, err)
	}
	
	
	helpers.CreateLog(user.UserID, "USER_REGISTER", user.UserID, "users")
	

	c.JSON(http.StatusOK, gin.H{"message": "Registrasi berhasil"})
}



// Login godoc
// @Summary Login user
// @Description Login dengan email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginInput true "Data login"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mohon lengkapi data login dan captcha"})
		return
	}

	if err := helpers.VerifyRecaptcha(input.RecaptchaToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verifikasi Captcha Gagal: Anda terdeteksi sebagai robot atau token kadaluwarsa."})
		return
	}
	

	var user models.User
    
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


	accessToken, err := helpers.CreateAccessToken(user.UserID, user.Role.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat access token"})
		return
	}

	
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Mendapatkan access token baru menggunakan refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object{refresh_token=string} true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token dibutuhkan"})
		return
	}

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


	userID := uint(claims["sub"].(float64))
	var user models.User
	

	if err := config.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if user.StatusUser == "inactive" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah dinonaktifkan."})
		return
	}


	newAccessToken, err := helpers.CreateAccessToken(user.UserID, user.Role.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat access token baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}


// RequestPasswordReset godoc
// @Summary Request reset password (authenticated)
// @Description Kirim kode verifikasi ke email untuk reset password
// @Tags Password & Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{current_password=string} true "Password saat ini"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /me/request-password-reset [post]
func RequestPasswordReset(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password saat ini harus diisi"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password saat ini tidak valid"})
		return
	}


	verificationCode := helpers.GenerateVerificationCode()
	expiresAt := time.Now().Add(10 * time.Minute) // Berlaku 10 menit


	config.DB.Where("user_id = ? AND is_used = false", user.UserID).Delete(&models.PasswordReset{})


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

// VerifyAndResetPassword godoc
// @Summary Verifikasi dan reset password
// @Description Memverifikasi kode OTP dan mengubah password
// @Tags Password & Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{verification_code=string,new_password=string} true "Kode verifikasi dan password baru"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /me/verify-reset-password [post]
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


// RequestEmailChange godoc
// @Summary Request ubah email
// @Description Kirim kode verifikasi ke email baru
// @Tags Password & Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{new_email=string} true "Email baru"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /me/request-email-change [post]
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

// VerifyAndChangeEmail godoc
// @Summary Verifikasi dan ubah email
// @Description Memverifikasi kode OTP dan mengubah email
// @Tags Password & Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{verification_code=string} true "Kode verifikasi"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /me/verify-email-change [post]
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

// ForgotPasswordRequest godoc
// @Summary Forgot password (tanpa login)
// @Description Kirim kode OTP ke email untuk reset password tanpa perlu login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object{email=string} true "Email terdaftar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /forgot-password [post]
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
	// fmt.Printf("[INFO] [FORGOT PASSWORD] Mengirim kode verifikasi ke email: %s\n", user.Email)
	// fmt.Printf("[INFO] Kode verifikasi: %s\n", verificationCode)
	
	if err := helpers.SendVerificationEmail(user.Email, user.Name, verificationCode, "reset_password"); err != nil {
		// fmt.Printf("[ERROR] Gagal mengirim email: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim email verifikasi"})
		return
	}
	
	// fmt.Printf("[SUCCESS] Email lupa password berhasil dikirim ke: %s\n", user.Email)

	helpers.CreateLog(user.UserID, "FORGOT_PASSWORD_REQUEST", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Kode verifikasi telah dikirim ke email Anda. Berlaku 10 menit.",
		"expires_at": expiresAt,
	})
}

// ForgotPasswordVerify godoc
// @Summary Verifikasi forgot password (tanpa login)
// @Description Verifikasi kode OTP dan reset password tanpa perlu login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object{email=string,verification_code=string,new_password=string} true "Data reset"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /forgot-password/verify [post]
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