package models

import "time"

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,max=16777215"`
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,max=255,min=8"`
}

type SignUpResponse struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
