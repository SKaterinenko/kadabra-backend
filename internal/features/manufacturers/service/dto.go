package manufacturers_service

type TranslationInput struct {
	LanguageCode string `json:"language_code"`
	Description  string `json:"description"`
}

type CreateInput struct {
	Name         string             `json:"name"`
	Translations []TranslationInput `json:"translations"`
}

type PatchInput struct {
	Name *string
}
