// controllers/user_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminCreateUser godoc
// @Summary Admin buat user baru
// @Description Membuat user baru dengan role spesifik (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{name=string,email=string,password=string,role_id=int} true "Data user baru"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /admin/users [post]
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

	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "CREATE_USER", user.UserID, "users")

	
	config.DB.Preload("Role").First(&user, user.UserID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat",
		"data":    user,
	})
}



// UpdateUser godoc
// @Summary Update data user
// @Description Mengubah data user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Param body body object{name=string,email=string} true "Data update"
// @Success 200 {object} map[string]interface{}
// @Router /admin/users/{id} [put]
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



// GetMyProfile godoc
// @Summary Profil saya
// @Description Mengambil data lengkap profil user yang login
// @Tags Profil
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /me [get]
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

// UpdateMyProfile godoc
// @Summary Update profil saya
// @Description Mengupdate profil user yang login (multipart/form-data)
// @Tags Profil
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string false "Nama"
// @Param bio formData string false "Bio"
// @Param image_file formData file false "Foto profil"
// @Success 200 {object} map[string]string
// @Router /me [put]
func UpdateMyProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	
	// Ambil data dari form
	name := c.PostForm("name")
	bio := c.PostForm("bio")
	
	// Update nama di tabel 'users' jika ada perubahan
	if name != "" {
		config.DB.Model(&user).Update("name", name)
	}

	// Update atau buat data di tabel 'user_profiles'
	var profile models.UserProfile
	config.DB.FirstOrInit(&profile, models.UserProfile{UserID: user.UserID})
	
	// Update bio 
	profile.Bio = bio
	
	// Debug log
	fmt.Printf("Updating profile for user %d: bio='%s'\n", user.UserID, bio)
	
	// Handle image upload ke ImageKit
	file, err := c.FormFile("image_file")
	if err == nil {
		// Buka file untuk upload
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca file"})
			return
		}
		defer fileContent.Close()

		// Upload ke ImageKit
		imageURL, err := helpers.UploadToImageKit(fileContent, file, "profiles")
		if err != nil {
			log.Printf("ImageKit upload error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar profil ke ImageKit"})
			return
		}

		// Update profile picture URL
		profile.ProfilePictureURL = imageURL
	}
	
	if err := config.DB.Save(&profile).Error; err != nil {
		fmt.Printf("Error saving profile: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil"})
		return
	}
	
	fmt.Printf("Profile saved successfully: bio='%s', picture='%s'\n", profile.Bio, profile.ProfilePictureURL)
	
	
	helpers.CreateLog(user.UserID, "UPDATE_PROFILE", user.UserID, "users")

	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui"})
}



// DeactivateUser godoc
// @Summary Nonaktifkan user
// @Description Mengubah status user menjadi inactive (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id}/deactivate [patch]
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

// ActivateUser godoc
// @Summary Aktifkan user
// @Description Mengubah status user menjadi active (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id}/activate [patch]
func ActivateUser(c *gin.Context) {
    var user models.User

    
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

    
	config.DB.Model(&user).Update("status_user", "active")
	admin := c.MustGet("user").(models.User)
	helpers.CreateLog(admin.UserID, "ADMIN_ACTIVATE_USER", user.UserID, "users")
	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil diaktifkan"})
}


// GetAllUsers godoc
// @Summary Ambil semua user
// @Description Mengambil semua user dengan informasi lengkap (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /admin/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User
	
	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUserByID godoc
// @Summary Detail user
// @Description Mengambil detail satu user (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Success 200 {object} map[string]interface{}
// @Router /admin/users/{id} [get]
func GetUserByID(c *gin.Context) {
	var user models.User
	
	if err := config.DB.Preload("Role").Preload("UserProfile").First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUserRole godoc
// @Summary Ubah role user
// @Description Mengubah role user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Param body body object{role_id=int} true "Role ID baru"
// @Success 200 {object} map[string]interface{}
// @Router /admin/users/{id}/role [patch]
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

	config.DB.Preload("Role").First(&user, user.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Role pengguna berhasil diubah",
		"data":    user,
	})
}

// GetUserRelatedData godoc
// @Summary Data terkait user
// @Description Mengambil jumlah data terkait user sebelum delete (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Success 200 {object} map[string]interface{}
// @Router /admin/users/{id}/related-data [get]
func GetUserRelatedData(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var menusCount int64
	var commentsCount int64
	var votesCount int64
	var bookmarksCount int64

	config.DB.Model(&models.Menu{}).Where("user_id = ?", userID).Count(&menusCount)
	config.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentsCount)
	config.DB.Model(&models.MenuVote{}).Where("user_id = ?", userID).Count(&votesCount)
	config.DB.Table("user_bookmarks").Where("user_id = ?", userID).Count(&bookmarksCount)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user_id":         user.UserID,
			"user_name":       user.Name,
			"menus_count":     menusCount,
			"comments_count":  commentsCount,
			"votes_count":     votesCount,
			"bookmarks_count": bookmarksCount,
			"has_related_data": menusCount > 0 || commentsCount > 0 || votesCount > 0 || bookmarksCount > 0,
		},
	})
}

// DeleteUser godoc
// @Summary Hapus user
// @Description Menghapus user beserta semua data terkaitnya (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id} [delete]
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

	// Cek apakah user yang akan dihapus adalah admin
	var targetUserRole models.Role
	config.DB.First(&targetUserRole, user.RoleID)
	if targetUserRole.RoleName == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak dapat menghapus user dengan role admin"})
		return
	}

	userIDtoLog := user.UserID
	userIDtoDelete := user.UserID

	tx := config.DB.Begin()

	// 1. Hapus user_bookmarks 
	if err := tx.Exec("DELETE FROM user_bookmarks WHERE user_id = ?", userIDtoDelete).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus bookmark pengguna"})
		return
	}

	// 2. Hapus menu_votes
	if err := tx.Where("user_id = ?", userIDtoDelete).Delete(&models.MenuVote{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus votes pengguna"})
		return
	}

	// 3. Hapus comments 
	
	if err := tx.Exec("DELETE FROM comments WHERE parent_id IN (SELECT comment_id FROM (SELECT comment_id FROM comments WHERE user_id = ?) AS subquery)", userIDtoDelete).Error; err != nil {
		
	}
	if err := tx.Where("user_id = ?", userIDtoDelete).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar pengguna"})
		return
	}

	// 4. Hapus menus beserta relasinya
	
	if err := tx.Exec("DELETE FROM menu_tags WHERE menu_id IN (SELECT menu_id FROM menus WHERE user_id = ?)", userIDtoDelete).Error; err != nil {
		
	}
	
	if err := tx.Exec("DELETE FROM comments WHERE menu_id IN (SELECT menu_id FROM menus WHERE user_id = ?)", userIDtoDelete).Error; err != nil {
		
	}
	
	if err := tx.Exec("DELETE FROM menu_votes WHERE menu_id IN (SELECT menu_id FROM menus WHERE user_id = ?)", userIDtoDelete).Error; err != nil {
		
	}
	
	if err := tx.Exec("DELETE FROM user_bookmarks WHERE menu_id IN (SELECT menu_id FROM menus WHERE user_id = ?)", userIDtoDelete).Error; err != nil {
		
	}
	// Hapus menus
	if err := tx.Where("user_id = ?", userIDtoDelete).Delete(&models.Menu{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus resep pengguna"})
		return
	}

	// 5. Hapus user_profile
	if err := tx.Where("user_id = ?", userIDtoDelete).Delete(&models.UserProfile{}).Error; err != nil {
		
	}

	// 6. Hapus notifications terkait user ini 
	if err := tx.Exec("DELETE FROM notifications WHERE user_id = ? OR actor_id = ?", userIDtoDelete, userIDtoDelete).Error; err != nil {
		
	}

	// 7. Hapus password_resets
	if err := tx.Exec("DELETE FROM password_resets WHERE user_id = ?", userIDtoDelete).Error; err != nil {
		
	}

	// 8. hapus user
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	// Commit transaksi
	tx.Commit()

	helpers.CreateLog(currentUser.UserID, "DELETE_USER", userIDtoLog, "users")

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna dan semua data terkait berhasil dihapus"})
}

// GetAllRoles godoc
// @Summary Ambil semua role
// @Description Mengambil semua role yang tersedia (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles [get]
func GetAllRoles(c *gin.Context) {
	var roles []models.Role
	
	if err := config.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// GetActivityLogs godoc
// @Summary Log aktivitas
// @Description Mengambil log aktivitas sistem dengan pagination (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman" default(1)
// @Param limit query int false "Jumlah per halaman" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /admin/logs [get]
func GetActivityLogs(c *gin.Context) {
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := strings.TrimSpace(c.Query("search"))
	dateFilter := c.Query("date_filter") 
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	
	var logs []models.LogActivity
	var total int64
	
	
	query := config.DB.Model(&models.LogActivity{})
	
	
	now := time.Now()
	switch dateFilter {
	case "today":
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		query = query.Where("created_at >= ?", startOfDay)
	case "week":
		query = query.Where("created_at >= ?", now.AddDate(0, 0, -7))
	case "month":
		query = query.Where("created_at >= ?", now.AddDate(0, -1, 0))
	case "year":
		query = query.Where("created_at >= ?", now.AddDate(-1, 0, 0))
	case "custom":
		if startDate != "" {
			query = query.Where("created_at >= ?", startDate)
		}
		if endDate != "" {
			query = query.Where("created_at <= ?", endDate+" 23:59:59")
		}
	}
	
	
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(action) LIKE ? OR LOWER(target_type) LIKE ? OR user_id IN (SELECT user_id FROM users WHERE LOWER(name) LIKE ?)",
			searchLower, searchLower, searchLower,
		)
	}
	
	
	query.Count(&total)
	
	
	if err := query.
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil log aktivitas"})
		return
	}
	
	
	totalPages := 1
	if total > 0 {
		totalPages = int(total) / limit
		if int(total)%limit > 0 {
			totalPages++
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"limit":        limit,
		},
	})
}