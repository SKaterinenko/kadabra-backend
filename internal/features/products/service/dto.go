package products_service

import ()

type CreateInput struct {
	Name             string
	Description      string
	ShortDescription string
	ProductsTypeId   int
	ManufacturerId   int
}

type PatchInput struct {
	Name             *string
	Description      *string
	ShortDescription *string
	ProductsTypeId   *int
	ManufacturerId   *int
}
