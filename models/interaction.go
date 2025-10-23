package models
// AddTagInput mendefinisikan input untuk menambahkan tag ke resep.
type AddTagInput struct {
	TagID uint `json:"tag_id" binding:"required"`
}