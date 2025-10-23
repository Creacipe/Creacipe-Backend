package models

// UserProfile menyimpan data tambahan pengguna.
type UserProfile struct {
	ProfileID         uint   `gorm:"primaryKey;column:profile_id"`
	UserID            uint   `gorm:"not null;unique"`
	Bio               string `gorm:"type:text"`
	ProfilePictureURL string `gorm:"type:varchar(255)"`
}