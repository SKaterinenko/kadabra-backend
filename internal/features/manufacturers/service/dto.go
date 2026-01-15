package manufacturers_service

type TranslationInput struct {
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type CreateInput struct {
	Translations []TranslationInput `json:"translations"`
}

type PatchInput struct {
	Name *string
}
