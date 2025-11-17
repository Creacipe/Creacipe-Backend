package models

import "time"

// PasswordReset menyimpan kode verifikasi untuk reset password
type PasswordReset struct {
	ResetID          uint      `gorm:"primaryKey;column:reset_id" json:"reset_id"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	VerificationCode string    `gorm:"size:6;not null" json:"-"` // 6 digit code
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	IsUsed           bool      `gorm:"default:false" json:"is_used"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// EmailVerification menyimpan kode verifikasi untuk ubah email
type EmailVerification struct {
	VerificationID   uint      `gorm:"primaryKey;column:verification_id" json:"verification_id"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	NewEmail         string    `gorm:"size:255;not null" json:"new_email"`
	VerificationCode string    `gorm:"size:6;not null" json:"-"` // 6 digit code
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	IsUsed           bool      `gorm:"default:false" json:"is_used"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}