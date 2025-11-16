package models

// Category mengelompokkan jenis-jenis tag.
type Category struct {
	CategoryID   uint   `gorm:"primaryKey;column:category_id" json:"category_id"`
	CategoryName string `gorm:"size:100;not null;column:category_name" json:"category_name"`
}

// CreateCategoryInput mendefinisikan input untuk membuat kategori baru.
type CreateCategoryInput struct {
	CategoryName string `json:"category_name" binding:"required"`
}

// UpdateCategoryInput mendefinisikan input untuk memperbarui kategori.
type UpdateCategoryInput struct {
	CategoryName string `json:"category_name" binding:"required"`
}