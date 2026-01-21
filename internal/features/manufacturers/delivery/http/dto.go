package manufacturers_http

type TranslationInput struct {
	LanguageCode string `json:"language_code" validate:"required"`
	Description  string `json:"description" validate:"required"`
}

type createDTO struct {
	Name         string             `json:"name" validate:"required"`
	Translations []TranslationInput `json:"translations" validate:"required,min=1,dive"`
}

type patchDTO struct {
	Translations *[]TranslationInput `json:"translations" db:"translations"`
	CategoryIds  *[]int              `json:"category_ids" db:"category_ids"`
}
