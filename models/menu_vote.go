package models

import "time"

// MenuVote menyimpan data like dan dislike dari user.
type MenuVote struct {
	VoteID        uint      `gorm:"primaryKey;column:vote_id"`
	MenuID        uint      `gorm:"not null"`
	UserID        uint      `gorm:"not null"`
	LikesCount    int       `gorm:"type:int;default:0;not null"` // 1 if user liked, 0 if not
	DislikesCount int       `gorm:"type:int;default:0;not null"` // 1 if user disliked, 0 if not
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// VoteInput untuk request body saat user vote
type VoteInput struct {
	VoteType string `json:"vote_type" binding:"required,oneof=like dislike"` // "like" or "dislike"
}