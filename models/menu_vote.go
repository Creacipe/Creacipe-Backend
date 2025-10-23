package models

import "time"

// MenuVote menyimpan data like (+1) atau dislike (-1) dari user.
type MenuVote struct {
	VoteID    uint      `gorm:"primaryKey;column:vote_id"`
	MenuID    uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	VoteType  int       `gorm:"type:tinyint;not null"` // 1 for like, -1 for dislike
	CreatedAt time.Time
	UpdatedAt time.Time
}