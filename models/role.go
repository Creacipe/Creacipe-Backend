package models

// Role mendefinisikan peran pengguna (admin, editor, member).
type Role struct {
	RoleID   uint   `gorm:"primaryKey;column:role_id" json:"role_id"`
	RoleName string `gorm:"size:50;not null;unique;column:role_name" json:"role_name"`
}