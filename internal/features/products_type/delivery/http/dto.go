package products_type_http

type createDTO struct {
	Translations  []TranslationInput `json:"translations" validate:"required,min=1,dive"`
	SubCategoryId int                `json:"sub_category_id" validate:"required"`
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}

type TranslationInput struct {
	LanguageCode string `json:"language_code" validate:"required"`
	Name         string `json:"name" validate:"required"`
}
