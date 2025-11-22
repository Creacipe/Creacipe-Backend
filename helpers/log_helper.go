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

// CreateNotification membuat notifikasi untuk user
func CreateNotification(userID uint, title string, message string, notifType string, relatedID *uint, relatedType string) {
	notification := models.Notification{
		UserID:      userID,
		Title:       title,
		Message:     message,
		Type:        notifType,
		IsRead:      false,
		RelatedID:   relatedID,
		RelatedType: relatedType,
		CreatedAt:   time.Now(),
	}
	config.DB.Create(&notification)
}