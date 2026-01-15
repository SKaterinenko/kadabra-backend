package categories_service

type TranslationInput struct {
	LanguageCode string `json:"language_code" binding:"required,len=2,alpha"`
	Name         string `json:"name" binding:"required,min=1,max=255"`
}

type CreateInput struct {
	Translations []TranslationInput `json:"translations" binding:"required,min=1,dive"`
}

type PatchInput struct {
	Name *string
}
