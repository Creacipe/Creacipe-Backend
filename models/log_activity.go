package models

import "time"

// LogActivity menyimpan jejak semua aktivitas penting pengguna.
type LogActivity struct {
	ActivityID uint      `gorm:"primaryKey;column:activity_id"`
	UserID     uint      `gorm:"not null"`
	Action     string    `gorm:"type:varchar(100);not null"`
	TargetID   uint
	TargetType string    `gorm:"type:varchar(100)"`
	CreatedAt  time.Time

	User User `gorm:"foreignKey:UserID"`
}

// TableName memberitahu GORM nama tabel yang benar di database.
func (LogActivity) TableName() string {
    return "log_activity"
}


