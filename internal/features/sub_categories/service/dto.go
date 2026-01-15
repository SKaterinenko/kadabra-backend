package sub_categories_service

type Translations []TranslationInput

type TranslationInput struct {
	LanguageCode string `json:"language_code" db:"language_code"`
	Name         string `json:"name" db:"name"`
}

type CreateInput struct {
	Translations Translations `json:"translations" validate:"required" db:"translations"`
	CategoryId   int          `json:"category_id" validate:"required" db:"category_id"`
}

type PatchInput struct {
	Name *string
}
