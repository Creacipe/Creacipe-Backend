package models

import "time"

// Menu adalah model utama untuk resep.
type Menu struct {
	MenuID          uint      `gorm:"primaryKey;column:menu_id" json:"menu_id"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	Title           string    `gorm:"size:255;not null" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	Ingredients     string    `gorm:"type:text" json:"ingredients"`
	Instructions    string    `gorm:"type:text" json:"instructions"`
	ImageURL        string    `gorm:"type:varchar(255)" json:"image_url"`
	Status          string    `gorm:"type:enum('pending','approved','rejected');default:'pending'" json:"status"`
	RejectionReason string    `gorm:"type:text" json:"rejection_reason"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// --- Relasi ---
	User     User       `gorm:"foreignKey:UserID;references:UserID" json:"User,omitempty"`
	Comments []Comment  `gorm:"foreignKey:MenuID" json:"-"`
	Votes    []MenuVote `gorm:"foreignKey:MenuID" json:"-"`
	Tags     []*Tag     `gorm:"many2many:menu_tags;foreignKey:MenuID;joinForeignKey:menu_id;References:TagID;joinReferences:tag_id" json:"tags,omitempty"`
}


type CreateMenuInput struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Ingredients  string `json:"ingredients" binding:"required"`
	Instructions string `json:"instructions" binding:"required"`
	ImageURL     string `json:"image_url"`
	TagIDs       []uint `json:"tag_ids"` 
}

type UpdateMenuInput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
	ImageURL     string `json:"image_url"`
}