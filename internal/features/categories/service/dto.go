package categories_service

type TranslationInput struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

type CreateInput struct {
	Translations map[string]TranslationInput `json:"translations" binding:"required"`
}

type PatchInput struct {
	Name *string
}
