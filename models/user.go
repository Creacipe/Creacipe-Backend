// models/user.go
package models

import "time"

type User struct {
	UserID    uint   `gorm:"primaryKey;column:user_id"`
	Name      string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255);unique;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Role      string `gorm:"type:enum('admin', 'editor', 'member');default:'member'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// --- DEFINISI RELASI YANG LEBIH LENGKAP ---
	Menus     []Menu       `gorm:"foreignKey:UserID"`
	Ratings   []MenuRating `gorm:"foreignKey:UserID"`
	Bookmarks []*Menu      `gorm:"many2many:user_bookmarks;foreignKey:UserID;joinForeignKey:user_id;References:MenuID;joinReferences:menu_id"`
}