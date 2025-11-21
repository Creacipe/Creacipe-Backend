// controllers/user_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers" // <-- 1. IMPORT HELPER
	"creacipe-backend/models"
	"net/http"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// --- FUNGSI BARU UNTUK ADMIN ---

// AdminCreateUser membuat user baru dengan peran spesifik.
func AdminCreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		RoleID   uint   `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi role exists
	var role models.Role
	if err := config.DB.First(&role, input.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		Name:       input.Name,
		Email:      input.Email,
		Password:   string(hashedPassword),
		RoleID:     input.RoleID,
		StatusUser: "active",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Buat profil kosong otomatis
	profile := models.UserProfile{UserID: user.UserID}
	config.DB.Create(&profile)

	// Log aktivitas
	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "CREATE_USER", user.UserID, "users")

	// Reload user dengan relasi
	config.DB.Preload("Role").First(&user, user.UserID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat",
		"data":    user,
	})
}



// UpdateUser menangani perubahan data pengguna (nama/email) oleh admin.
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password" binding:"omitempty,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pengguna"})
		return
	}

	// Log aktivitas
	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "UPDATE_USER", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil diperbarui",
		"data":    user,
	})
}



//--------------------------------------// USER PROFILE CONTROLLERS //--------------------------------------//

// GetMyProfile menampilkan data lengkap dari pengguna yang sedang login.
func GetMyProfile(c *gin.Context) {
	userCtx, _ := c.Get("user")
	user := userCtx.(models.User)

	var userDetails models.User
	if err := config.DB.Preload("Role").Preload("Profile").First(&userDetails, user.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": userDetails})
}

// UpdateMyProfile memperbarui data profil pengguna yang sedang login.
func UpdateMyProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	
	// Ambil data dari form
	name := c.PostForm("name")
	bio := c.PostForm("bio")
	
	// Update nama di tabel 'users' jika ada perubahan
	// Note: Email TIDAK bisa diubah di sini, gunakan endpoint change-email
	if name != "" {
		config.DB.Model(&user).Update("name", name)
	}

	// Update atau buat data di tabel 'user_profiles'
	var profile models.UserProfile
	config.DB.FirstOrInit(&profile, models.UserProfile{UserID: user.UserID})
	
	// Update bio - pastikan update meskipun kosong
	profile.Bio = bio
	
	// Debug log
	fmt.Printf("Updating profile for user %d: bio='%s'\n", user.UserID, bio)
	
	// Handle image upload
	file, err := c.FormFile("image_file")
	if err == nil {
		// Buat nama file unik
		filename := fmt.Sprintf("profile_%d_%d%s", user.UserID, time.Now().Unix(), filepath.Ext(file.Filename))
		filePath := "assets/profiles/" + filename
		
		// Simpan file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar profil"})
			return
		}
		
		// Update profile picture URL
		profile.ProfilePictureURL = "/" + filePath
	}
	
	if err := config.DB.Save(&profile).Error; err != nil {
		fmt.Printf("Error saving profile: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil"})
		return
	}
	
	fmt.Printf("Profile saved successfully: bio='%s', picture='%s'\n", profile.Bio, profile.ProfilePictureURL)
	
	// Catat di log
	helpers.CreateLog(user.UserID, "UPDATE_PROFILE", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui"})
}

//--------------------------------------// ADMIN USER MANAGEMENT //--------------------------------------//

// --- FUNGSI BARU UNTUK AKTIF/NONAKTIF ---

// DeactivateUser mengubah status user menjadi 'inactive'.
func DeactivateUser(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}
	
	config.DB.Model(&user).Update("status_user", "inactive")
	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "ADMIN_DEACTIVATE_USER", user.UserID, "users")
	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dinonaktifkan"})
}

// ActivateUser mengubah status user menjadi 'active'.
func ActivateUser(c *gin.Context) {
    var user models.User

    // --- PERBAIKAN DI SINI ---
    // Gunakan GORM secara langsung, ia akan otomatis mencari berdasarkan primary key 'user_id'
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}
    // -------------------------
    
	config.DB.Model(&user).Update("status_user", "active")
	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "ADMIN_ACTIVATE_USER", user.UserID, "users")
	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil diaktifkan"})
}

// --- FUNGSI TAMBAHAN UNTUK DASHBOARD ADMIN ---

// GetAllUsers mengambil semua pengguna dengan informasi lengkap
func GetAllUsers(c *gin.Context) {
	var users []models.User
	
	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUserByID mengambil detail satu pengguna
func GetUserByID(c *gin.Context) {
	var user models.User
	
	if err := config.DB.Preload("Role").Preload("UserProfile").First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUserRole mengubah role user (Admin only)
func UpdateUserRole(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var input struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID wajib diisi"})
		return
	}

	// Validasi role exists
	var role models.Role
	if err := config.DB.First(&role, input.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	user.RoleID = input.RoleID

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah role pengguna"})
		return
	}

	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "UPDATE_USER_ROLE", user.UserID, "users")

	// Reload user dengan role
	config.DB.Preload("Role").First(&user, user.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Role pengguna berhasil diubah",
		"data":    user,
	})
}

// DeleteUser menghapus user (Admin only)
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Cek apakah user mencoba menghapus dirinya sendiri
	currentUser := c.MustGet("user").(models.User)
	if currentUser.UserID == user.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak dapat menghapus akun sendiri"})
		return
	}

	userIDtoLog := user.UserID

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	helpers.CreateLog(currentUser.UserID, "DELETE_USER", userIDtoLog, "users")

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}

// GetAllRoles mengambil semua role yang tersedia
func GetAllRoles(c *gin.Context) {
	var roles []models.Role
	
	if err := config.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// GetActivityLogs mengambil log aktivitas sistem (Admin only)
func GetActivityLogs(c *gin.Context) {
	var logs []models.LogActivity
	
	// Preload User dan Role terlebih dahulu
	if err := config.DB.
		Preload("User.Role").
		Preload("User").
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil log aktivitas"})
		return
	}

	// Preload Menu, TargetUser, Tag, dan Category secara conditional berdasarkan target_type
	for i := range logs {
		targetType := logs[i].TargetType
		
		if targetType == "menu" || targetType == "menus" {
			// Preload Menu untuk menu-related actions
			config.DB.Preload("Menu").First(&logs[i], logs[i].ActivityID)
		} else if targetType == "user" || targetType == "users" {
			// Preload TargetUser untuk user-related actions
			config.DB.Preload("TargetUser.Role").Preload("TargetUser").First(&logs[i], logs[i].ActivityID)
		} else if targetType == "tag" || targetType == "tags" {
			// Preload Tag untuk tag-related actions
			config.DB.Preload("Tag").First(&logs[i], logs[i].ActivityID)
		} else if targetType == "category" || targetType == "categories" {
			// Preload Category untuk category-related actions
			config.DB.Preload("Category").First(&logs[i], logs[i].ActivityID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}