package tests

import (
	"bytes"
	"creacipe-backend/controllers"
	"creacipe-backend/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// ============ REGISTER TESTS ============

func TestRegister_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/register", controllers.Register)

	input := map[string]interface{}{
		"name":            "Test User",
		"email":           "test@example.com",
		"password":        "password123",
		"recaptcha_token": "test-token",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Debug: Print response body if failed
	if resp.Code != http.StatusOK {
		t.Logf("Response Body: %s", resp.Body.String())
	}

	assert.Equal(t, http.StatusOK, resp.Code, "Register should return 200 OK")
	
	var user models.User
	err := db.Preload("Role").Where("email = ?", "test@example.com").First(&user).Error
	assert.NoError(t, err, "User should be created in database")
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, uint(3), user.RoleID) // Default role member
}

func TestRegister_InvalidInput(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()
	router.POST("/register", controllers.Register)

	input := map[string]interface{}{
		"email": "invalid-email",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/register", controllers.Register)

	// Create existing user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := models.User{
		Name:     "Existing User",
		Email:    "duplicate@example.com",
		Password: string(hashedPassword),
		RoleID:   3,
	}
	db.Create(&existingUser)

	// Try to register with same email
	input := map[string]interface{}{
		"name":            "New User",
		"email":           "duplicate@example.com",
		"password":        "password123",
		"recaptcha_token": "test-token",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Email sudah terdaftar")
}

// ============ LOGIN TESTS ============

func TestLogin_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/login", controllers.Login)

	// Create user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      "login@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "active",
	}
	db.Create(&user)

	input := map[string]interface{}{
		"email":           "login@example.com",
		"password":        "password123",
		"recaptcha_token": "test-token",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_token"])
}

func TestLogin_InvalidEmail(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()
	router.POST("/login", controllers.Login)

	input := map[string]interface{}{
		"email":           "nonexistent@example.com",
		"password":        "password123",
		"recaptcha_token": "test-token",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Email atau password salah")
}

func TestLogin_InvalidPassword(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/login", controllers.Login)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      "test@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "active",
	}
	db.Create(&user)

	input := map[string]interface{}{
		"email":           "test@example.com",
		"password":        "wrongpassword",
		"recaptcha_token": "test-token",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Email atau password salah")
}

func TestLogin_InactiveUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/login", controllers.Login)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Inactive User",
		Email:      "inactive@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "inactive",
	}
	db.Create(&user)

	input := map[string]interface{}{
		"email":           "inactive@example.com",
		"password":        "password123",
		"recaptcha_token": "test-token",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Contains(t, resp.Body.String(), "Akun Anda telah dinonaktifkan")
}

// ============ REFRESH TOKEN TESTS ============

func TestRefreshToken_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/refresh", controllers.RefreshToken)

	// Create user with role
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      "refresh@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "active",
	}
	db.Create(&user)

	// Login to get refresh token
	router.POST("/login", controllers.Login)
	loginInput := map[string]interface{}{
		"email":           "refresh@example.com",
		"password":        "password123",
		"recaptcha_token": "test-token",
	}
	loginBody, _ := json.Marshal(loginInput)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)

	var loginResponse map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginResponse)
	refreshToken := loginResponse["refresh_token"].(string)

	// Use refresh token
	input := map[string]interface{}{
		"refresh_token": refreshToken,
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NotEmpty(t, response["access_token"])
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()
	router.POST("/refresh", controllers.RefreshToken)

	input := map[string]interface{}{
		"refresh_token": "invalid.token.here",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestRefreshToken_MissingToken(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()
	router.POST("/refresh", controllers.RefreshToken)

	input := map[string]interface{}{}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Controller returns 401 for invalid/missing token, not 400
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// ============ FORGOT PASSWORD TESTS ============

func TestForgotPasswordRequest_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/forgot-password", controllers.ForgotPasswordRequest)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      "forgot@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "active",
	}
	db.Create(&user)

	input := map[string]interface{}{
		"email": "forgot@example.com",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	
	// Verify password reset record created
	var passwordReset models.PasswordReset
	err := db.Where("user_id = ?", user.UserID).First(&passwordReset).Error
	assert.NoError(t, err)
	assert.False(t, passwordReset.IsUsed)
	assert.Len(t, passwordReset.VerificationCode, 6)
}

func TestForgotPasswordRequest_NonexistentEmail(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()
	router.POST("/forgot-password", controllers.ForgotPasswordRequest)

	input := map[string]interface{}{
		"email": "nonexistent@example.com",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Should still return 200 for security (don't reveal if email exists)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestForgotPasswordVerify_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/forgot-password/verify", controllers.ForgotPasswordVerify)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      "verify@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "active",
	}
	db.Create(&user)

	// Create password reset record
	passwordReset := models.PasswordReset{
		UserID:           user.UserID,
		VerificationCode: "123456",
		ExpiresAt:        time.Now().Add(10 * time.Minute),
		IsUsed:           false,
	}
	db.Create(&passwordReset)

	input := map[string]interface{}{
		"email":             "verify@example.com",
		"verification_code": "123456",
		"new_password":      "newpassword123",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/forgot-password/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	
	// Verify password was changed
	var updatedUser models.User
	db.First(&updatedUser, user.UserID)
	err := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword123"))
	assert.NoError(t, err)
	
	// Verify reset code is marked as used
	var usedReset models.PasswordReset
	db.Where("user_id = ?", user.UserID).First(&usedReset)
	assert.True(t, usedReset.IsUsed)
}

func TestForgotPasswordVerify_InvalidCode(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/forgot-password/verify", controllers.ForgotPasswordVerify)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      "test@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "active",
	}
	db.Create(&user)

	input := map[string]interface{}{
		"email":             "test@example.com",
		"verification_code": "999999",
		"new_password":      "newpassword123",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/forgot-password/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Kode verifikasi tidak valid")
}

func TestForgotPasswordVerify_ExpiredCode(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	router.POST("/forgot-password/verify", controllers.ForgotPasswordVerify)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:       "Test User",
		Email:      "expired@example.com",
		Password:   string(hashedPassword),
		RoleID:     3,
		StatusUser: "active",
	}
	db.Create(&user)

	// Create expired password reset
	passwordReset := models.PasswordReset{
		UserID:           user.UserID,
		VerificationCode: "123456",
		ExpiresAt:        time.Now().Add(-1 * time.Minute), // Expired 1 minute ago
		IsUsed:           false,
	}
	db.Create(&passwordReset)

	input := map[string]interface{}{
		"email":             "expired@example.com",
		"verification_code": "123456",
		"new_password":      "newpassword123",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/forgot-password/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "kadaluarsa")
}

// ============ REQUEST PASSWORD RESET (ME) TESTS ============

func TestRequestPasswordReset_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "resetme@example.com", 3)

	router.POST("/me/request-password-reset", createAuthContext(user), controllers.RequestPasswordReset)

	input := map[string]interface{}{
		"current_password": "password123",
	}
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/me/request-password-reset", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// ============ VERIFY AND RESET PASSWORD (ME) TESTS ============

func TestVerifyAndResetPassword_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "verifyme@example.com", 3)

	// Create password reset
	passwordReset := models.PasswordReset{
		UserID:           user.UserID,
		VerificationCode: "123456",
		ExpiresAt:        time.Now().Add(10 * time.Minute),
		IsUsed:           false,
	}
	db.Create(&passwordReset)

	router.POST("/me/verify-reset-password", createAuthContext(user), controllers.VerifyAndResetPassword)

	input := map[string]interface{}{
		"verification_code": "123456",
		"new_password":      "newpassword123",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/me/verify-reset-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// ============ REQUEST EMAIL CHANGE TESTS ============

func TestRequestEmailChange_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "olduser@example.com", 3)

	router.POST("/me/request-email-change", createAuthContext(user), controllers.RequestEmailChange)

	input := map[string]interface{}{
		"new_email": "newuser@example.com",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/me/request-email-change", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// ============ VERIFY AND CHANGE EMAIL TESTS ============

func TestVerifyAndChangeEmail_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "oldchange@example.com", 3)

	// Create email verification request using SQL directly
	db.Exec("INSERT INTO email_verifications (user_id, new_email, verification_code, expires_at, is_used) VALUES (?, ?, ?, ?, ?)",
		user.UserID, "newchange@example.com", "654321", time.Now().Add(10*time.Minute), false)

	router.POST("/me/verify-email-change", createAuthContext(user), controllers.VerifyAndChangeEmail)

	input := map[string]interface{}{
		"verification_code": "654321",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/me/verify-email-change", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
