package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllMenusForModeration menampilkan semua resep (semua status) untuk editor/admin.
func GetAllMenusForModeration(c *gin.Context) {
    var menus []models.Menu
    config.DB.Preload("User").Find(&menus)
    c.JSON(http.StatusOK, gin.H{"data": menus})
}

// GetPendingMenus menampilkan resep yang menunggu persetujuan ('pending').
func GetPendingMenus(c *gin.Context) {
    var menus []models.Menu
    config.DB.Preload("User").Where("status = ?", "pending").Find(&menus)
    c.JSON(http.StatusOK, gin.H{"data": menus})
}

// UpdateMenuStatusInput mendefinisikan input untuk mengubah status resep.
type UpdateMenuStatusInput struct {
    Status          string `json:"status" binding:"required,oneof=approved,rejected"`
    RejectionReason string `json:"rejection_reason"`
}

// UpdateMenuStatus menangani logika approve/reject resep oleh editor/admin.
func UpdateMenuStatus(c *gin.Context) {
    var menu models.Menu
    if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
        return
    }

    var input UpdateMenuStatusInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updates := map[string]interface{}{"status": input.Status}
    if input.Status == "rejected" {
        updates["rejection_reason"] = input.RejectionReason
    } else if input.Status == "approved" {
        updates["rejection_reason"] = ""
    }

    config.DB.Model(&menu).Updates(updates)
    c.JSON(http.StatusOK, gin.H{"data": menu})
}