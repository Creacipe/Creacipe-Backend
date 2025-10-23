package models

import "time"

// Comment menyimpan data komentar pada resep.
type Comment struct {
	CommentID   uint      `gorm:"primaryKey;column:comment_id"`
	MenuID      uint      `gorm:"not null"`
	UserID      uint      `gorm:"not null"`
	CommentText string    `gorm:"type:text;not null"`
	CreatedAt   time.Time
}