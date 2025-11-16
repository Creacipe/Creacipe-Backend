package models

import "time"

// LogActivity menyimpan jejak semua aktivitas penting pengguna.
type LogActivity struct {
	ActivityID uint      `gorm:"primaryKey;column:activity_id" json:"activity_id"`
	UserID     uint      `gorm:"not null;column:user_id" json:"user_id"`
	Action     string    `gorm:"type:varchar(100);not null" json:"action"`
	TargetID   uint      `gorm:"column:target_id" json:"target_id"`
	TargetType string    `gorm:"type:varchar(100);column:target_type" json:"target_type"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`

	User       User     `gorm:"foreignKey:UserID;references:UserID" json:"User"`
	TargetUser *User    `gorm:"foreignKey:TargetID;references:UserID" json:"TargetUser,omitempty"`
	Menu       *Menu    `gorm:"foreignKey:TargetID;references:MenuID" json:"Menu,omitempty"`
	Tag        *Tag     `gorm:"foreignKey:TargetID;references:TagID" json:"Tag,omitempty"`
	Category   *Category `gorm:"foreignKey:TargetID;references:CategoryID" json:"Category,omitempty"`
}

// TableName memberitahu GORM nama tabel yang benar di database.
func (LogActivity) TableName() string {
    return "log_activity"
}


