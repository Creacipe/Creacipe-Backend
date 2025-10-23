package models

import "time"

type PasswordReset struct {
	ResetID   uint      `gorm:"primaryKey;column:reset_id"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"size:255;not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
}