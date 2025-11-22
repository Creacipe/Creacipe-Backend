package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMyNotifications mengambil semua notifikasi user yang login
func GetMyNotifications(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var notifications []models.Notification
	query := config.DB.Where("user_id = ?", user.UserID).
		Order("created_at DESC")

	// Filter by read/unread (optional)
	if readFilter := c.Query("is_read"); readFilter != "" {
		isRead := readFilter == "true"
		query = query.Where("is_read = ?", isRead)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var total int64
	query.Model(&models.Notification{}).Count(&total)

	if err := query.Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetUnreadNotificationCount mengambil jumlah notifikasi yang belum dibaca
func GetUnreadNotificationCount(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var count int64
	if err := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", user.UserID).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// MarkNotificationAsRead menandai satu notifikasi sebagai sudah dibaca
func MarkNotificationAsRead(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	notifID := c.Param("id")

	var notification models.Notification
	if err := config.DB.Where("notification_id = ? AND user_id = ?", notifID, user.UserID).
		First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan"})
		return
	}

	if err := config.DB.Model(&notification).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi ditandai sudah dibaca"})
}

// MarkAllNotificationsAsRead menandai semua notifikasi user sebagai sudah dibaca
func MarkAllNotificationsAsRead(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	if err := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", user.UserID).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semua notifikasi ditandai sudah dibaca"})
}
