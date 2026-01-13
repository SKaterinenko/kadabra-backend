package products_type_service

import ()

type CreateInput struct {
	Name          string
	SubCategoryId int
}

type PatchInput struct {
	Name *string
}
