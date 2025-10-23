package models

// Category mengelompokkan jenis-jenis tag.
type Category struct {
	CategoryID   uint   `gorm:"primaryKey;column:category_id"`
	CategoryName string `gorm:"size:100;not null;column:category_name"`
}