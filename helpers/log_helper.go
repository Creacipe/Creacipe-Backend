package helpers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"time"
)

// CreateLog adalah fungsi terpusat untuk mencatat aktivitas pengguna.
func CreateLog(userID uint, action string, targetID uint, targetType string) {
	activity := models.LogActivity{
		UserID:     userID,
		Action:     action,
		TargetID:   targetID,
		TargetType: targetType,
		CreatedAt:  time.Now(),
	}
	config.DB.Create(&activity)
}