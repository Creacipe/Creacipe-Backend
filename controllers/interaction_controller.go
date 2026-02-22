// controllers/interaction_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddTagToMenu godoc
// @Summary Tambahkan tag ke resep
// @Description Menambahkan tag ke sebuah resep
// @Tags Interaksi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Param body body models.AddTagInput true "Tag ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/{id}/tags [post]
func AddTagToMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	var input models.AddTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}
	var tag models.Tag
	if err := config.DB.First(&tag, input.TagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag tidak ditemukan"})
		return
	}
	
	if err := config.DB.Model(&menu).Association("Tags").Append(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan tag"})
		return
	}

	// Catat aktivitas
	user := c.MustGet("user").(models.User)
	helpers.CreateLog(user.UserID, "ADD_TAG_TO_MENU", menu.MenuID, "menus")

	c.JSON(http.StatusOK, gin.H{"message": "Tag berhasil ditambahkan ke resep"})
}

// BookmarkMenu godoc
// @Summary Bookmark resep
// @Description Menambahkan resep ke daftar bookmark user
// @Tags Interaksi
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]string
// @Router /menus/{id}/bookmark [post]
func BookmarkMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	user := c.MustGet("user").(models.User)
	
	if err := config.DB.Model(&user).Association("Bookmarks").Append(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan bookmark"})
		return
	}

	// Catat aktivitas
	helpers.CreateLog(user.UserID, "BOOKMARK_MENU", menu.MenuID, "menus")

	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil di-bookmark"})
}

// UnbookmarkMenu godoc
// @Summary Hapus bookmark
// @Description Menghapus resep dari daftar bookmark user
// @Tags Interaksi
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]string
// @Router /menus/{id}/bookmark [delete]
func UnbookmarkMenu(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	user := c.MustGet("user").(models.User)
	
	if err := config.DB.Model(&user).Association("Bookmarks").Delete(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus bookmark"})
		return
	}

	// Catat aktivitas
	helpers.CreateLog(user.UserID, "UNBOOKMARK_MENU", menu.MenuID, "menus")

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark berhasil dihapus"})
}

// handleVote adalah fungsi internal untuk mengelola logika like/dislike.
func handleVote(c *gin.Context, voteType string) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	user := c.MustGet("user").(models.User)

	var vote models.MenuVote
	err := config.DB.Where("user_id = ? AND menu_id = ?", user.UserID, menu.MenuID).First(&vote).Error

	var actionLog string
	var message string

	if err == nil {
		// Jika sudah ada record vote sebelumnya
		if voteType == "like" {
			if vote.LikesCount == 1 {
				// User sudah like, sekarang unlike (toggle off)
				vote.LikesCount = 0
				actionLog = "UNLIKE_MENU"
				message = "Like berhasil dibatalkan"
			} else {
				// User belum like atau sedang dislike, sekarang like
				vote.LikesCount = 1
				vote.DislikesCount = 0 // Hapus dislike jika ada
				actionLog = "LIKE_MENU"
				message = "Resep berhasil di-like"
			}
		} else { // voteType == "dislike"
			if vote.DislikesCount == 1 {
				// User sudah dislike, sekarang un-dislike (toggle off)
				vote.DislikesCount = 0
				actionLog = "UNDISLIKE_MENU"
				message = "Dislike berhasil dibatalkan"
			} else {
				// User belum dislike atau sedang like, sekarang dislike
				vote.DislikesCount = 1
				vote.LikesCount = 0 // Hapus like jika ada
				actionLog = "DISLIKE_MENU"
				message = "Resep berhasil di-dislike"
			}
		}
		config.DB.Save(&vote)
	} else {
		// Jika belum ada record, buat baru
		newVote := models.MenuVote{
			UserID:        user.UserID,
			MenuID:        menu.MenuID,
			LikesCount:    0,
			DislikesCount: 0,
		}
		
		if voteType == "like" {
			newVote.LikesCount = 1
			actionLog = "LIKE_MENU"
			message = "Resep berhasil di-like"
		} else {
			newVote.DislikesCount = 1
			actionLog = "DISLIKE_MENU"
			message = "Resep berhasil di-dislike"
		}
		
		config.DB.Create(&newVote)
	}

	// Catat aktivitas
	helpers.CreateLog(user.UserID, actionLog, menu.MenuID, "menus")
	
	// Return current state
	var currentVote models.MenuVote
	config.DB.Where("user_id = ? AND menu_id = ?", user.UserID, menu.MenuID).First(&currentVote)
	
	c.JSON(http.StatusOK, gin.H{
		"message":        message,
		"is_liked":       currentVote.LikesCount == 1,
		"is_disliked":    currentVote.DislikesCount == 1,
	})
}

// LikeMenu godoc
// @Summary Like resep
// @Description Like/unlike resep (toggle)
// @Tags Interaksi
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]interface{}
// @Router /menus/{id}/like [post]
func LikeMenu(c *gin.Context) {
	handleVote(c, "like")
}

// DislikeMenu godoc
// @Summary Dislike resep
// @Description Dislike/un-dislike resep (toggle)
// @Tags Interaksi
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]interface{}
// @Router /menus/{id}/dislike [post]
func DislikeMenu(c *gin.Context) {
	handleVote(c, "dislike")
}

// GetUserInteractionStatus godoc
// @Summary Status interaksi
// @Description Cek status like/dislike/bookmark user terhadap resep
// @Tags Interaksi
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID resep"
// @Success 200 {object} map[string]interface{}
// @Router /menus/{id}/interaction-status [get]
func GetUserInteractionStatus(c *gin.Context) {
	var menu models.Menu
	if err := config.DB.First(&menu, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}
	
	user := c.MustGet("user").(models.User)
	
	// Cek status vote
	var vote models.MenuVote
	isLiked := false
	isDisliked := false
	err := config.DB.Where("user_id = ? AND menu_id = ?", user.UserID, menu.MenuID).First(&vote).Error
	if err == nil {
		isLiked = vote.LikesCount == 1
		isDisliked = vote.DislikesCount == 1
	}
	
	// Cek status bookmark
	var bookmarkCount int64
	config.DB.Table("user_bookmarks").Where("user_id = ? AND menu_id = ?", user.UserID, menu.MenuID).Count(&bookmarkCount)
	isBookmarked := bookmarkCount > 0
	
	c.JSON(http.StatusOK, gin.H{
		"is_liked":      isLiked,
		"is_disliked":   isDisliked,
		"is_bookmarked": isBookmarked,
	})
}