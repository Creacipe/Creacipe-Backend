package models

// UpdateMenuStatusInput mendefinisikan input JSON untuk mengubah status resep.
type UpdateMenuStatusInput struct {
	Status          string `json:"status" binding:"required,oneof=approved,rejected"`
	RejectionReason string `json:"rejection_reason"`
}