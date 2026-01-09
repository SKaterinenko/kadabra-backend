package categories_model

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewCategory(name, slug string) *Category {
	return &Category{
		Id:   uuid.New(),
		Name: name,
		Slug: slug,
	}
}

type CategoryPatch struct {
	Name *string `json:"name,omitempty" db:"name"`
}

func NewCategoryPatch(name string) *CategoryPatch {
	return &CategoryPatch{
		Name: &name,
	}
}
