// models/tag.go
package models

import "time"

type Tag struct {
	TagID     uint   `gorm:"primaryKey;column:tag_id"`
	Name      string `gorm:"type:varchar(100);unique;not null"`
	Type      string `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}