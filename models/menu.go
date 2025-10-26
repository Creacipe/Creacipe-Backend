package models

import "time"

// Menu adalah model utama untuk resep.
type Menu struct {
	MenuID          uint      `gorm:"primaryKey;column:menu_id"`
	UserID          uint      `gorm:"not null"`
	Title           string    `gorm:"size:255;not null"`
	Description     string    `gorm:"type:text"`
	Ingredients     string    `gorm:"type:text"`
	Instructions    string    `gorm:"type:text"`
	ImageURL        string    `gorm:"type:varchar(255)"`
	Status          string    `gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	RejectionReason string    `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// --- Relasi ---
	User     User       `gorm:"foreignKey:UserID"`
	Comments []Comment  `gorm:"foreignKey:MenuID"`
	Votes    []MenuVote `gorm:"foreignKey:MenuID"`
	Tags     []*Tag     `gorm:"many2many:menu_tags;foreignKey:MenuID;joinForeignKey:menu_id;References:TagID;joinReferences:tag_id"`
}

// --- Struct Input didefinisikan di sini ---
type CreateMenuInput struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Ingredients  string `json:"ingredients" binding:"required"`
	Instructions string `json:"instructions" binding:"required"`
	ImageURL     string `json:"image_url"`
	TagIDs       []uint `json:"tag_ids"` // TAMBAHKAN INI
}

type UpdateMenuInput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
	ImageURL     string `json:"image_url"`
}