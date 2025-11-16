package models

// UserProfile menyimpan data tambahan pengguna.
type UserProfile struct {
	ProfileID         uint   `gorm:"primaryKey;column:profile_id" json:"profile_id"`
	UserID            uint   `gorm:"not null;unique" json:"user_id"`
	Bio               string `gorm:"type:text" json:"bio"`
	ProfilePictureURL string `gorm:"type:varchar(255)" json:"profile_picture_url"`
}