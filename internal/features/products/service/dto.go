package products_service

import (
	"github.com/google/uuid"
)

type CreateInput struct {
	Name             string
	Description      string
	ShortDescription string
	ProductsTypeId   uuid.UUID
	ManufacturerId   uuid.UUID
}

type PatchInput struct {
	Name             *string
	Description      *string
	ShortDescription *string
	ProductsTypeId   *uuid.UUID
	ManufacturerId   *uuid.UUID
}
