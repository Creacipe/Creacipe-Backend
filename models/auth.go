package models


// RegisterInput mendefinisikan data JSON yang dibutuhkan untuk registrasi.
type RegisterInput struct {
	Name     string `json:"name" bson:"name" binding:"required"`
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required,min=6"`
}

// LoginInput mendefinisikan data JSON yang dibutuhkan untuk login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}