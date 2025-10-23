package models

// Tag menyimpan data tag yang bisa ditempelkan ke resep.
type Tag struct {
	TagID      uint   `gorm:"primaryKey;column:tag_id"`
	CategoryID uint   `gorm:"not null"`
	TagName    string `gorm:"size:100;not null;unique;column:tag_name"`

	Category Category `gorm:"foreignKey:CategoryID"`
}

// CreateTagInput mendefinisikan input JSON untuk membuat tag baru.
type CreateTagInput struct {
	TagName    string `json:"tag_name" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

// UpdateTagInput mendefinisikan input JSON untuk memperbarui tag.
type UpdateTagInput struct {
	TagName    string `json:"tag_name"`
	CategoryID uint   `json:"category_id"`
}