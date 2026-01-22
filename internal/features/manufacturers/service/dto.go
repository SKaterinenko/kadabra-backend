package manufacturers_service

type TranslationInput struct {
	LanguageCode string `json:"language_code"`
	Description  string `json:"description"`
}

type CreateInput struct {
	Name         string             `json:"name"`
	Translations []TranslationInput `json:"translations"`
	CategoryIds  []int              `json:"category_ids" db:"category_ids"`
}

type PatchInput struct {
	Translations *[]TranslationInput `json:"translations" db:"translations"`
	CategoryIds  *[]int              `json:"category_ids" db:"category_ids"`
}
