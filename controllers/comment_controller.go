package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment - User membuat komentar di resep
func CreateComment(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	menuID := c.Param("id")

	// Validasi: resep harus ada
	var menu models.Menu
	if err := config.DB.First(&menu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	// Bind input
	var input struct {
		CommentText string `json:"comment_text" binding:"required"`
		ParentID    *uint  `json:"parent_id"` // Optional, untuk reply
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Komentar tidak boleh kosong"})
		return
	}

	// Jika ini adalah reply, validasi parent comment exists
	if input.ParentID != nil {
		var parentComment models.Comment
		if err := config.DB.First(&parentComment, *input.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Komentar yang direply tidak ditemukan"})
			return
		}
		// Pastikan parent comment adalah untuk menu yang sama
		if parentComment.MenuID != menu.MenuID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Komentar tidak valid"})
			return
		}
	}

	// Buat comment
	comment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user.UserID,
		ParentID:    input.ParentID,
		CommentText: input.CommentText,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat komentar"})
		return
	}

	// Auto-create notification
	var notificationUserID uint
	var notificationMessage string
	
	if input.ParentID != nil {
		// Jika reply, notif ke pemilik comment parent
		var parentComment models.Comment
		config.DB.First(&parentComment, *input.ParentID)
		notificationUserID = parentComment.UserID
		notificationMessage = user.Name + " membalas komentar Anda: \"" + input.CommentText + "\""
		
		// Jangan buat notif jika reply ke diri sendiri
		if notificationUserID != user.UserID {
			notification := models.Notification{
				UserID:      notificationUserID,
				Title:       "Balasan Baru",
				Message:     notificationMessage,
				Type:        "info",
				IsRead:      false,
				RelatedID:   &menu.MenuID,
				RelatedType: "menu",
			}
			config.DB.Create(&notification)
		}
	} else if menu.UserID != user.UserID {
		// Jika comment biasa DAN bukan pemilik resep, notif ke pemilik resep
		notificationUserID = menu.UserID
		notificationMessage = user.Name + " berkomentar di resep Anda: \"" + input.CommentText + "\""
		
		notification := models.Notification{
			UserID:      notificationUserID,
			Title:       "Komentar Baru",
			Message:     notificationMessage,
			Type:        "info",
			IsRead:      false,
			RelatedID:   &menu.MenuID,
			RelatedType: "menu",
		}
		config.DB.Create(&notification)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Komentar berhasil ditambahkan",
		"data":    comment,
	})
}

// GetCommentsByMenu - Get semua komentar untuk suatu resep (dengan nested replies)
func GetCommentsByMenu(c *gin.Context) {
	menuID := c.Param("id")

	// Cek apakah kolom parent_id ada (untuk backward compatibility)
	var comments []models.Comment
	var err error
	
	// Coba query dengan parent_id filter
	err = config.DB.Where("menu_id = ? AND (parent_id IS NULL OR parent_id = 0)", menuID).
		Order("created_at DESC").
		Find(&comments).Error
	
	// Jika error karena kolom tidak ada, fallback ke query tanpa parent_id
	if err != nil {
		err = config.DB.Where("menu_id = ?", menuID).
			Order("created_at DESC").
			Find(&comments).Error
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil komentar"})
		return
	}

	var buildComment func(models.Comment) models.CommentWithUser
	buildComment = func(comment models.Comment) models.CommentWithUser {
		var user models.User
		config.DB.Preload("Profile").First(&user, comment.UserID)

		avatar := ""
		if user.Profile.ProfilePictureURL != "" {
			avatar = user.Profile.ProfilePictureURL
		}

		// Get replies untuk comment ini (jika kolom parent_id ada)
		var replies []models.Comment
		var repliesWithUser []models.CommentWithUser
		
		// Try to get replies, ignore error jika kolom belum ada
		err := config.DB.Where("parent_id = ?", comment.CommentID).
			Order("created_at ASC").
			Find(&replies).Error
		
		if err == nil {
			for _, reply := range replies {
				repliesWithUser = append(repliesWithUser, buildComment(reply))
			}
		}

		return models.CommentWithUser{
			CommentID:   comment.CommentID,
			MenuID:      comment.MenuID,
			UserID:      comment.UserID,
			ParentID:    comment.ParentID,
			UserName:    user.Name,
			UserAvatar:  avatar,
			CommentText: comment.CommentText,
			CreatedAt:   comment.CreatedAt.Format("2006-01-02 15:04:05"),
			Replies:     repliesWithUser,
		}
	}

	var result []models.CommentWithUser
	for _, comment := range comments {
		result = append(result, buildComment(comment))
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
	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	// Cek apakah user adalah pemilik comment
	if comment.UserID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus komentar ini"})
		return
	}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dihapus"})
}
