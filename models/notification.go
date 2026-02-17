package models

import "time"

// Notification menyimpan notifikasi untuk user
type Notification struct {
	NotificationID uint      `gorm:"primaryKey;column:notification_id" json:"notification_id"`
	UserID         uint      `gorm:"not null;column:user_id" json:"user_id"`
	Title          string    `gorm:"type:varchar(255);not null" json:"title"`
	Message        string    `gorm:"type:text;not null" json:"message"`
	Type           string    `gorm:"type:enum('info','success','warning','danger');default:'info'" json:"type"`
	IsRead         bool      `gorm:"default:false;column:is_read" json:"is_read"`
	RelatedID      *uint     `gorm:"column:related_id" json:"related_id,omitempty"`
	RelatedType    string    `gorm:"type:varchar(50);column:related_type" json:"related_type,omitempty"` 
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`

	User User  `gorm:"foreignKey:UserID;references:UserID" json:"User,omitempty"`
	Menu *Menu `gorm:"foreignKey:RelatedID;references:MenuID" json:"Menu,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
