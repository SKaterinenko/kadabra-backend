package sub_categories_service

import ()

type CreateInput struct {
	Name       string
	CategoryId int
}

type PatchInput struct {
	Name *string
}
