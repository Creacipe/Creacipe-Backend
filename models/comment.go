package models

import "time"


// Comment menyimpan data komentar pada resep.
type Comment struct {
	CommentID   uint      `gorm:"primaryKey;column:comment_id" json:"comment_id"`
	MenuID      uint      `gorm:"not null" json:"menu_id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	ParentID    *uint     `gorm:"default:null" json:"parent_id"` // Untuk reply
	CommentText string    `gorm:"type:text;not null" json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Relasi untuk nested replies (optional, bisa di-load saat perlu)
	Replies     []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

type CommentWithUser struct {
		CommentID   uint              `json:"comment_id"`
		MenuID      uint              `json:"menu_id"`
		UserID      uint              `json:"user_id"`
		ParentID    *uint             `json:"parent_id"`
		UserName    string            `json:"user_name"`
		UserAvatar  string            `json:"user_avatar"`
		CommentText string            `json:"comment_text"`
		CreatedAt   string            `json:"created_at"`
		Replies     []CommentWithUser `json:"replies"`
	}