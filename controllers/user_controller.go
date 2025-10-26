// controllers/user_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers" // <-- 1. IMPORT HELPER
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// --- FUNGSI BARU UNTUK ADMIN ---

// AdminCreateUser membuat user baru dengan peran spesifik.
func AdminCreateUser(c *gin.Context) {
	var input models.AdminCreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role models.Role
	if err := config.DB.Where("role_name = ?", input.RoleName).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama peran tidak valid"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		Name:       input.Name,
		Email:      input.Email,
		Password:   string(hashedPassword),
		RoleID:     role.RoleID,
		StatusUser: "active",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// --- TAMBAHAN: BUAT PROFIL KOSONG OTOMATIS ---
	profile := models.UserProfile{UserID: user.UserID}
	config.DB.Create(&profile)
	// ---------------------------------------------
	// --- 2. TAMBAHKAN LOG UNTUK PENAMBAHAN USER ---

	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "ADMIN_CREATE_USER", user.UserID, "users")
	c.JSON(http.StatusCreated, gin.H{"data": user})
}


// GetAllUsers menampilkan semua pengguna, bisa difilter berdasarkan status.
func GetAllUsers(c *gin.Context) {
	var users []models.User
	query := config.DB.Preload("Role")

	status := c.Query("status")
	if status == "inactive" {
		query = query.Where("status_user = ?", "inactive")
	} else if status == "all" {
		// Tidak melakukan apa-apa, ambil semua status
	} else {
		// Defaultnya hanya ambil user yang aktif
		query = query.Where("status_user = ?", "active")
	}
	
	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}



// UpdateUser menangani perubahan data pengguna (nama/email) oleh admin.
func UpdateUser(c *gin.Context) {
	var user models.User // User yang akan diubah
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var input models.AdminUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil data admin yang melakukan aksi dari context
	admin, _ := c.Get("user")
	adminInfo := admin.(models.User)

	if err := config.DB.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pengguna"})
		return
	}

	// --- 2. TAMBAHKAN LOG ---
	helpers.CreateLog(adminInfo.UserID, "ADMIN_UPDATE_USER", user.UserID, "users")
	// -------------------------

	c.JSON(http.StatusOK, gin.H{"data": user})
}



// UpdateUserRole berfungsi untuk mengubah peran seorang pengguna.
func UpdateUserRole(c *gin.Context) {
	var user models.User // User yang akan diubah perannya
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var input models.UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil data admin yang melakukan aksi dari context
	admin, _ := c.Get("user")
	adminInfo := admin.(models.User)

	var role models.Role
	if err := config.DB.Where("role_name = ?", input.RoleName).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama peran tidak valid"})
		return
	}

	if err := config.DB.Model(&user).Update("role_id", role.RoleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui peran pengguna"})
		return
	}

	// --- 3. TAMBAHKAN LOG ---
	helpers.CreateLog(adminInfo.UserID, "ADMIN_UPDATE_ROLE", user.UserID, "users")
	// -------------------------

	config.DB.Preload("Role").First(&user, user.UserID)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser menangani penghapusan pengguna oleh admin.
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Ambil data admin yang melakukan aksi dari context
	admin, _ := c.Get("user")
	adminInfo := admin.(models.User)

	// Simpan ID user yang akan dihapus untuk log
	userIDtoLog := user.UserID

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	// --- 4. TAMBAHKAN LOG ---
	helpers.CreateLog(adminInfo.UserID, "ADMIN_DELETE_USER", userIDtoLog, "users")
	// -------------------------

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
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
	var input models.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update nama di tabel 'users'
	config.DB.Model(&user).Update("name", input.Name)

	// Update atau buat data di tabel 'user_profiles'
	var profile models.UserProfile
	config.DB.FirstOrInit(&profile, models.UserProfile{UserID: user.UserID})
	profile.Bio = input.Bio
	profile.ProfilePictureURL = input.ProfilePictureURL
	config.DB.Save(&profile)
	
	// Catat di log
	helpers.CreateLog(user.UserID, "UPDATE_PROFILE", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui"})
}

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