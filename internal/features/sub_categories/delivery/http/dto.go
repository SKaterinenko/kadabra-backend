package sub_categories_http

import "github.com/google/uuid"

type createDTO struct {
	Name       string    `json:"name" validate:"required"`
	CategoryId uuid.UUID `json:"categoryId" validate:"required"`
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}
