// controllers/tag_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllTags menampilkan semua tag yang tersedia di database.
func GetAllTags(c *gin.Context) {
	var tags []models.Tag
	config.DB.Find(&tags)
	c.JSON(http.StatusOK, gin.H{"data": tags})
}