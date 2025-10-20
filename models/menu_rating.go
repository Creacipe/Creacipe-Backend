// models/menu_rating.go
package models

import "time"

type MenuRating struct {
	RatingID  uint   `gorm:"primaryKey;column:rating_id"`
	UserID    uint   `gorm:"column:user_id;not null"`
	MenuID    uint   `gorm:"column:menu_id;not null"`
	Rating    uint   `gorm:"type:tinyint;not null"` // Rating 1-5
	CreatedAt time.Time
	UpdatedAt time.Time
}