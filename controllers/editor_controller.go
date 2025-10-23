// controllers/editor_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers" // <-- 1. IMPORT HELPER
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllMenusForModeration menampilkan semua resep (semua status) untuk editor/admin.
func GetAllMenusForModeration(c *gin.Context) {
	var menus []models.Menu
	if err := config.DB.Preload("User").Order("created_at desc").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data resep"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// GetPendingMenus menampilkan resep yang menunggu persetujuan ('pending').
func GetPendingMenus(c *gin.Context) {
	var menus []models.Menu
	if err := config.DB.Preload("User").Where("status = ?", "pending").Order("created_at asc").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data resep pending"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// UpdateMenuStatus menangani logika approve/reject resep oleh editor/admin.
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

	// --- 2. TAMBAHKAN LOG UNTUK AKSI MODERASI ---
	actionLog := "APPROVE_MENU"
	if input.Status == "rejected" {
		actionLog = "REJECT_MENU"
	}
	helpers.CreateLog(moderatorInfo.UserID, actionLog, menu.MenuID, "menus")
	// ---------------------------------------------
	
	c.JSON(http.StatusOK, gin.H{"data": menu})
}