package sub_categories_http

type createDTO struct {
	Name       string `json:"name" validate:"required"`
	CategoryId int    `json:"categoryId" validate:"required"`
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}
