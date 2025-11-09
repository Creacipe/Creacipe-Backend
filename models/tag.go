package models

// Tag menyimpan data tag yang bisa ditempelkan ke resep.
type Tag struct {
	TagID      uint   `gorm:"primaryKey;column:tag_id" json:"tag_id"`
	CategoryID uint   `gorm:"not null" json:"category_id"`
	TagName    string `gorm:"size:100;not null;unique;column:tag_name" json:"tag_name"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
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