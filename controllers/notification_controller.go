package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMyNotifications godoc
// @Summary Ambil notifikasi
// @Description Mengambil notifikasi user yang login dengan pagination
// @Tags Notifikasi
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman" default(1)
// @Param limit query int false "Jumlah per halaman" default(20)
// @Param is_read query string false "Filter: true/false"
// @Success 200 {object} map[string]interface{}
// @Router /me/notifications [get]
func GetMyNotifications(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var notifications []models.Notification
	query := config.DB.Where("user_id = ?", user.UserID).
		Order("created_at DESC")

	if readFilter := c.Query("is_read"); readFilter != "" {
		isRead := readFilter == "true"
		query = query.Where("is_read = ?", isRead)
	}

	
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

// GetUnreadNotificationCount godoc
// @Summary Hitung notifikasi belum dibaca
// @Description Mengambil jumlah notifikasi yang belum dibaca
// @Tags Notifikasi
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int
// @Router /me/notifications/unread-count [get]
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

// MarkNotificationAsRead godoc
// @Summary Tandai notifikasi dibaca
// @Description Menandai satu notifikasi sebagai sudah dibaca
// @Tags Notifikasi
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID notifikasi"
// @Success 200 {object} map[string]string
// @Router /me/notifications/{id}/read [patch]
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

// MarkAllNotificationsAsRead godoc
// @Summary Tandai semua notifikasi dibaca
// @Description Menandai semua notifikasi user sebagai sudah dibaca
// @Tags Notifikasi
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /me/notifications/mark-all-read [patch]
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
