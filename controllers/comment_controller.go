package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment godoc
// @Summary Tambah komentar
// @Description Membuat komentar di resep. Bisa reply komentar lain dengan parent_id
// @Tags Komentar
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Param body body object{comment_text=string,parent_id=int} true "Data komentar"
// @Success 201 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /menus/{id}/comments [post]
func CreateComment(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	menuID := c.Param("id")

	// 1. Validasi: Resep harus ada
	var menu models.Menu
	if err := config.DB.First(&menu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	// 2. Bind input JSON
	var input struct {
		CommentText string `json:"comment_text" binding:"required"`
		ParentID    *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Komentar tidak boleh kosong"})
		return
	}

	var finalParentID *uint     
	var notificationUserID uint 

	if input.ParentID != nil {
		var targetComment models.Comment
		
		if err := config.DB.First(&targetComment, *input.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Komentar yang direply tidak ditemukan"})
			return
		}

		if targetComment.MenuID != menu.MenuID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Komentar tidak valid untuk resep ini"})
			return
		}

		notificationUserID = targetComment.UserID


		if targetComment.ParentID != nil {
			finalParentID = targetComment.ParentID
		} else {
			
			finalParentID = input.ParentID
		}

	} else {
	
		finalParentID = nil
	
		if menu.UserID != user.UserID {
			notificationUserID = menu.UserID
		}
	}

	// 4. Simpan Komentar ke Database
	comment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user.UserID,
		ParentID:    finalParentID, 
		CommentText: input.CommentText,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat komentar"})
		return
	}

	
	if notificationUserID != 0 && notificationUserID != user.UserID {
		var notifMessage string
		var notifTitle string

		if input.ParentID != nil {
			notifTitle = "Balasan Baru"
			notifMessage = user.Name + " membalas komentar Anda: \"" + input.CommentText + "\""
		} else {
			notifTitle = "Komentar Baru"
			notifMessage = user.Name + " berkomentar di resep Anda: \"" + input.CommentText + "\""
		}

		notification := models.Notification{
			UserID:      notificationUserID,
			Title:       notifTitle,
			Message:     notifMessage,
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

// GetCommentsByMenu godoc
// @Summary Ambil komentar
// @Description Mengambil semua komentar untuk suatu resep (termasuk replies nested)
// @Tags Komentar
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]interface{}
// @Router /menus/{id}/comments [get]
func GetCommentsByMenu(c *gin.Context) {
	menuID := c.Param("id")

	
	var comments []models.Comment
	var err error
	
	
	err = config.DB.Where("menu_id = ? AND (parent_id IS NULL OR parent_id = 0)", menuID).
		Order("created_at DESC").
		Find(&comments).Error
	
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

		
		var replies []models.Comment
		var repliesWithUser []models.CommentWithUser
		
		
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

// DeleteComment godoc
// @Summary Hapus komentar
// @Description Menghapus komentar milik sendiri
// @Tags Komentar
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID komentar"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	commentID, _ := strconv.Atoi(c.Param("id"))

	var comment models.Comment

	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	if comment.UserID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses untuk menghapus komentar ini"})
		return
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar dari database"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dihapus"})
}
