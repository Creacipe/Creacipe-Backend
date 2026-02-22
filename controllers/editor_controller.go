// controllers/editor_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats godoc
// @Summary Statistik dashboard
// @Description Mengambil statistik untuk dashboard editor/admin
// @Tags Editor
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /editor/dashboard/stats [get]
func GetDashboardStats(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	// Hitung resep pending
	var pendingCount int64
	config.DB.Model(&models.Menu{}).Where("status = ?", "pending").Count(&pendingCount)

	// Hitung total resep approved
	var approvedCount int64
	config.DB.Model(&models.Menu{}).Where("status = ?", "approved").Count(&approvedCount)

	// Hitung total resep rejected
	var rejectedCount int64
	config.DB.Model(&models.Menu{}).Where("status = ?", "rejected").Count(&rejectedCount)

	// Hitung total pengguna (khusus admin)
	var totalUsers int64
	if user.Role.RoleName == "admin" {
		config.DB.Model(&models.User{}).Count(&totalUsers)
	}

	// Hitung total kategori
	var totalCategories int64
	config.DB.Model(&models.Category{}).Count(&totalCategories)

	// Hitung total tag
	var totalTags int64
	config.DB.Model(&models.Tag{}).Count(&totalTags)

	stats := gin.H{
		"pending_recipes":  pendingCount,
		"approved_recipes": approvedCount,
		"rejected_recipes": rejectedCount,
		"total_categories": totalCategories,
		"total_tags":       totalTags,
	}

	// Tambahkan data user untuk admin
	if user.Role.RoleName == "admin" {
		stats["total_users"] = totalUsers
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetAllMenusForModeration godoc
// @Summary Semua resep (moderasi)
// @Description Menampilkan semua resep untuk moderasi (semua status)
// @Tags Editor
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter: pending, approved, rejected"
// @Success 200 {object} map[string]interface{}
// @Router /editor/menus [get]
func GetAllMenusForModeration(c *gin.Context) {
	var menus []models.Menu
	
	// Ambil parameter status dari query string
	status := c.Query("status")
	
	query := config.DB.Preload("User").Preload("User.Role")
	
	// Filter berdasarkan status jika ada
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Order("created_at desc").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data resep"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// GetPendingMenus godoc
// @Summary Resep pending
// @Description Menampilkan resep yang menunggu persetujuan
// @Tags Editor
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /editor/menus/pending [get]
func GetPendingMenus(c *gin.Context) {
	var menus []models.Menu
	if err := config.DB.Preload("User").Preload("User.Role").Where("status = ?", "pending").Order("created_at asc").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data resep pending"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// UpdateMenuStatus godoc
// @Summary Approve/reject resep
// @Description Mengubah status resep (approve/reject) oleh editor/admin
// @Tags Editor
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Param body body models.UpdateMenuStatusInput true "Status baru"
// @Success 200 {object} map[string]interface{}
// @Router /editor/menus/{id}/status [patch]
func UpdateMenuStatus(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	var input models.UpdateMenuStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil data editor/admin yang melakukan aksi
	moderator, _ := c.Get("user")
	moderatorInfo := moderator.(models.User)

	updates := map[string]interface{}{"status": input.Status}
	if input.Status == "rejected" {
		updates["rejection_reason"] = input.RejectionReason
	} else if input.Status == "approved" {
		updates["rejection_reason"] = ""
	}

	if err := config.DB.Model(&menu).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status resep"})
		return
	}

	actionLog := "APPROVE_MENU"
	if input.Status == "rejected" {
		actionLog = "REJECT_MENU"
	}
	helpers.CreateLog(moderatorInfo.UserID, actionLog, menu.MenuID, "menus")

	if input.Status == "approved" {
		helpers.CreateNotification(
			menu.UserID,
			"Resep Disetujui 🎉",
			"Resep \""+menu.Title+"\" telah disetujui dan dipublikasikan!",
			"success",
			&menu.MenuID,
			"menu",
		)
	} else if input.Status == "rejected" {
		reason := input.RejectionReason
		if reason == "" {
			reason = "Tidak ada alasan yang diberikan"
		}
		helpers.CreateNotification(
			menu.UserID,
			"Resep Ditolak",
			"Resep \""+menu.Title+"\" ditolak. Alasan: "+reason,
			"danger",
			&menu.MenuID,
			"menu",
		)
	}
	// -----------------------------------------------
	
	c.JSON(http.StatusOK, gin.H{"data": menu})
}