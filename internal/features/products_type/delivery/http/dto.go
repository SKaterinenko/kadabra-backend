package products_type_http

import "github.com/google/uuid"

type createDTO struct {
	Name          string    `json:"name" validate:"required"`
	SubCategoryId uuid.UUID `json:"subCategoryId" validate:"required"`
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}
