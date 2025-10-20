// models/menu.go
package models

import "time"

type Menu struct {
	MenuID       uint   `gorm:"primaryKey;column:menu_id"`
	UserID       uint   `gorm:"column:user_id;not null"`
	Title        string `gorm:"type:varchar(255);not null"`
	Description  string `gorm:"type:text"`
	Ingredients  string `gorm:"type:json"`
	Instructions string `gorm:"type:text"`
	ImageURL     string `gorm:"type:varchar(255)"`
	Status       string `gorm:"type:enum('pending', 'approved', 'rejected');default:'pending'"` 
	RejectionReason string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// --- DEFINISI RELASI YANG LEBIH LENGKAP ---
	User      User         `gorm:"foreignKey:UserID"`
	Tags      []*Tag       `gorm:"many2many:menu_tags;foreignKey:MenuID;joinForeignKey:menu_id;References:TagID;joinReferences:tag_id"`
	Ratings   []MenuRating `gorm:"foreignKey:MenuID"`
}