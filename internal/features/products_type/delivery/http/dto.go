package products_type_http

type createDTO struct {
	Name          string `json:"name" validate:"required"`
	SubCategoryId int    `json:"subCategoryId" validate:"required"`
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}
