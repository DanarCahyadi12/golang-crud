package models

import "go-crud/internal/entity"

type ProductRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Price int    `json:"price" validate:"min=0"`
	Stock int    `json:"stock" validate:"min=0"`
}

type ProductResponse struct {
	Id    string        `json:"id,omitempty"`
	Name  string        `json:"name,omitempty"`
	Price int           `json:"price,omitempty"`
	Stock int           `json:"stock,omitempty"`
	User  []entity.User `json:"user,omitempty"`
}
