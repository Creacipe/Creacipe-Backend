package models

import "time"

// User mendefinisikan data utama pengguna untuk login dan relasi.
type User struct {
	UserID    uint   `gorm:"primaryKey;column:user_id"`
	RoleID    uint   `gorm:"not null"`
	StatusUser string `gorm:"type:enum('active','inactive');default:'active';column:status_user"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;not null;unique"`
	Password  string `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// --- Relasi ---
	Role      Role
	Profile   UserProfile `gorm:"foreignKey:UserID"`
	Menus     []Menu      `gorm:"foreignKey:UserID"`
	Comments  []Comment   `gorm:"foreignKey:UserID"`
	Votes     []MenuVote  `gorm:"foreignKey:UserID"`
	Bookmarks []*Menu     `gorm:"many2many:user_bookmarks;foreignKey:UserID;joinForeignKey:user_id;References:MenuID;joinReferences:menu_id"`
}

// AdminCreateUserInput mendefinisikan input untuk membuat user baru oleh admin.
type AdminCreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleName string `json:"role_name" binding:"required,oneof=admin,editor,member"`
}

// AdminUpdateUserInput mendefinisikan data user yang bisa diubah oleh admin.
type AdminUpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}

// UpdateUserRoleInput mendefinisikan input saat mengubah peran user.
type UpdateUserRoleInput struct {
	RoleName string `json:"role_name" binding:"required,oneof=admin editor member"`
}

// UpdateProfileInput mendefinisikan data JSON untuk memperbarui profil.
type UpdateProfileInput struct {
	Name              string `json:"name"`
	Bio               string `json:"bio"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

