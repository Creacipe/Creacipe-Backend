package models

import "time"

// User mendefinisikan data utama pengguna untuk login dan relasi.
type User struct {
	UserID    uint   `gorm:"primaryKey;column:user_id" json:"user_id"`
	RoleID    uint   `gorm:"not null" json:"role_id"`
	StatusUser string `gorm:"type:enum('active','inactive');default:'active';column:status_user" json:"status_user"`
	Name      string `gorm:"size:255;not null" json:"name"`
	Email     string `gorm:"size:255;not null;unique" json:"email"`
	Password  string `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// --- Relasi ---
	Role      Role 		    `gorm:"foreignKey:RoleID;references:RoleID" json:"Role"`
	Profile   UserProfile `gorm:"foreignKey:UserID" json:"Profile"`
	UserProfile UserProfile `gorm:"foreignKey:UserID" json:"-"` //alias untuk code coverage
	Menus     []Menu      `gorm:"foreignKey:UserID" json:"-"`
	Comments  []Comment   `gorm:"foreignKey:UserID" json:"-"`
	Votes     []MenuVote  `gorm:"foreignKey:UserID" json:"-"`
	Bookmarks []*Menu     `gorm:"many2many:user_bookmarks;foreignKey:UserID;joinForeignKey:user_id;References:MenuID;joinReferences:menu_id" json:"-"`
}

// AdminCreateUserInput mendefinisikan input untuk membuat user baru oleh admin.
type AdminCreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	// PERBAIKAN: Tambahkan spasi di antara pilihan oneof
	RoleName string `json:"role_name" binding:"required,oneof=admin editor member"`
}

// AdminUpdateUserInput mendefinisikan data user yang bisa diubah oleh admin.
type AdminUpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}

// UpdateUserRoleInput mendefinisikan input saat mengubah peran user.
type UpdateUserRoleInput struct {
	// PERBAIKAN: Tambahkan spasi di antara pilihan oneof
	RoleName string `json:"role_name" binding:"required,oneof=admin editor member"`
}

// UpdateProfileInput mendefinisikan data JSON untuk memperbarui profil.
type UpdateProfileInput struct {
	Name              string `json:"name"`
	Email             string `json:"email" binding:"omitempty,email"`
	Bio               string `json:"bio"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

