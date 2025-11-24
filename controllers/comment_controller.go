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

	// Variabel penampung logika
	var finalParentID *uint     
	var notificationUserID uint 

	// 3. Logika Penentuan Parent & Notifikasi
	if input.ParentID != nil {
		// --- KASUS REPLY (BALASAN) ---
		var targetComment models.Comment
		
		// Cari komentar yang sedang dibalas
		if err := config.DB.First(&targetComment, *input.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Komentar yang direply tidak ditemukan"})
			return
		}

		// Validasi: Pastikan komentar target ada di menu yang sama
		if targetComment.MenuID != menu.MenuID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Komentar tidak valid untuk resep ini"})
			return
		}
		// A. Tentukan Penerima Notifikasi
		notificationUserID = targetComment.UserID

		// B. Tentukan Struktur Database (Flattening):
		// Jika komentar yang dibalas SUDAH punya parent (artinya dia adalah reply),
		// maka komentar baru kita harus nempel ke Parent-nya dia (Root).
		if targetComment.ParentID != nil {
			finalParentID = targetComment.ParentID
		} else {
			// Jika komentar yang dibalas adalah Root (ParentID nil),
			// maka kita nempel langsung ke dia.
			finalParentID = input.ParentID
		}

	} else {
		// --- KASUS KOMENTAR BARU (ROOT) ---
		finalParentID = nil
		// Notifikasi dikirim ke Pemilik Resep
		if menu.UserID != user.UserID {
			notificationUserID = menu.UserID
		}
	}

	// 4. Simpan Komentar ke Database
	comment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user.UserID,
		ParentID:    finalParentID, // Gunakan ID yang sudah diluruskan logic-nya
		CommentText: input.CommentText,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat komentar"})
		return
	}

	// 5. Kirim Notifikasi (Jika perlu)
	// Hanya kirim jika penerima bukan diri sendiri
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

// DeleteComment
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

	// 4. Berikan respon sukses
	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dihapus"})
}
