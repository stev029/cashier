package models

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty,min=6,max=20"`
	GroupID  []uint `json:"group_id" binding:"required"`
	Role     string `json:"role"`
}
