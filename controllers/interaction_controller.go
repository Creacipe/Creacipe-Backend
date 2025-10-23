// controllers/interaction_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/helpers" // <-- IMPORT HELPER
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddTagToMenu menangani logika penambahan tag ke sebuah resep.
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

// BookmarkMenu menangani penambahan resep ke daftar bookmark pengguna.
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

// UnbookmarkMenu menangani penghapusan resep dari daftar bookmark.
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
func handleVote(c *gin.Context, voteType int) {
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

	if err == nil { // Jika sudah ada vote sebelumnya
		if vote.VoteType == voteType {
			// Jika user menekan tombol yang sama, hapus vote-nya (toggle off)
			config.DB.Delete(&vote)
			actionLog = "REMOVE_VOTE"
			message = "Vote berhasil ditarik"
		} else {
			// Jika user mengubah vote (misal dari dislike ke like), update
			vote.VoteType = voteType
			config.DB.Save(&vote)
			actionLog = "UPDATE_VOTE"
			message = "Vote berhasil diubah"
		}
	} else { // Jika belum ada vote, buat baru
		newVote := models.MenuVote{
			UserID:   user.UserID,
			MenuID:   menu.MenuID,
			VoteType: voteType,
		}
		config.DB.Create(&newVote)
		actionLog = "CREATE_VOTE"
		message = "Vote berhasil disimpan"
	}
	
	// Tentukan log spesifik untuk like/dislike
	if actionLog == "CREATE_VOTE" || actionLog == "UPDATE_VOTE" {
		if voteType == 1 { actionLog = "LIKE_MENU" } else { actionLog = "DISLIKE_MENU" }
	}

	// Catat aktivitas
	helpers.CreateLog(user.UserID, actionLog, menu.MenuID, "menus")
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// LikeMenu menangani saat user menekan tombol 'like'.
func LikeMenu(c *gin.Context) {
	handleVote(c, 1) // 1 artinya 'like'
}

// DislikeMenu menangani saat user menekan tombol 'dislike'.
func DislikeMenu(c *gin.Context) {
	handleVote(c, -1) // -1 artinya 'dislike'
}