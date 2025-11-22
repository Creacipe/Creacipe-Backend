package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment - User membuat komentar di resep (tidak bisa comment di resep sendiri)
func CreateComment(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	menuID := c.Param("id")

	// Validasi: user tidak bisa comment di resep sendiri
	var menu models.Menu
	if err := config.DB.First(&menu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	if menu.UserID == user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak bisa berkomentar di resep sendiri"})
		return
	}

	// Bind input
	var input struct {
		CommentText string `json:"comment_text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Komentar tidak boleh kosong"})
		return
	}

	// Buat comment
	comment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user.UserID,
		CommentText: input.CommentText,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat komentar"})
		return
	}

	// Auto-create notification untuk pemilik resep
	notification := models.Notification{
		UserID:      menu.UserID, // Pemilik resep
		Title:       "Komentar Baru",
		Message:     user.Name + " berkomentar: \"" + input.CommentText + "\"",
		Type:        "info",
		IsRead:      false,
		RelatedID:   &menu.MenuID,
		RelatedType: "menu",
	}
	config.DB.Create(&notification)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Komentar berhasil ditambahkan",
		"data":    comment,
	})
}

// GetCommentsByMenu - Get semua komentar untuk suatu resep
func GetCommentsByMenu(c *gin.Context) {
	menuID := c.Param("id")

	var comments []models.Comment
	if err := config.DB.Where("menu_id = ?", menuID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil komentar"})
		return
	}

	

	var result []models.CommentWithUser
	for _, comment := range comments {
		var user models.User
		config.DB.Preload("Profile").First(&user, comment.UserID)

		avatar := ""
		if user.Profile.ProfilePictureURL != "" {
			avatar = user.Profile.ProfilePictureURL
		}

		result = append(result, models.CommentWithUser{
			CommentID:   comment.CommentID,
			MenuID:      comment.MenuID,
			UserID:      comment.UserID,
			UserName:    user.Name,
			UserAvatar:  avatar,
			CommentText: comment.CommentText,
			CreatedAt:   comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": len(result),
	})
}

// DeleteComment - User hapus komentar sendiri
func DeleteComment(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	commentID, _ := strconv.Atoi(c.Param("id"))

	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	// Validasi: hanya bisa hapus komentar sendiri
	if comment.UserID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak bisa menghapus komentar orang lain"})
		return
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dihapus"})
}
