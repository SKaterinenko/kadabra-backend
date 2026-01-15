package products_type_service

type TranslationInput struct {
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
}

type CreateInput struct {
	Translations  []TranslationInput `json:"translations"`
	SubCategoryId int                `json:"sub_category_id"`
}

type PatchInput struct {
	Name *string
}
