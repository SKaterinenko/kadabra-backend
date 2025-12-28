package model

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCategory(name string) *Category {
	return &Category{
		Id:   uuid.New(),
		Name: name,
	}
}

type CategoryPatch struct {
	Name *string `json:"name,omitempty"`
}

func NewCategoryPatch(name string) *CategoryPatch {
	return &CategoryPatch{
		Name: &name,
	}
}
