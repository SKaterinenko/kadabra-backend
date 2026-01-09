package products_http

import "github.com/google/uuid"

type createDTO struct {
	Name             string    `json:"name" validate:"required"`
	Description      string    `json:"description" validate:"required"`
	ShortDescription string    `json:"short_description" validate:"required"`
	ProductsTypeId   uuid.UUID `json:"products_type_id" validate:"required"`
	ManufacturerId   uuid.UUID `json:"manufacturer_id" validate:"required"`
}

type patchDTO struct {
	Name             *string    `json:"name,omitempty" validate:"omitempty"`
	Description      *string    `json:"description,omitempty" validate:"omitempty"`
	ShortDescription *string    `json:"short_description,omitempty" validate:"omitempty"`
	ProductsTypeId   *uuid.UUID `json:"products_type_id,omitempty" validate:"omitempty"`
	ManufacturerId   *uuid.UUID `json:"manufacturer_id,omitempty" validate:"omitempty"`
}

type getByIdsDTO struct {
	Ids []uuid.UUID `json:"ids" validate:"required"`
}
