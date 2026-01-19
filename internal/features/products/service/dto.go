package products_service

type TranslationInput struct {
	LanguageCode     string `json:"language_code"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
}

type CreateInput struct {
	Translations   []TranslationInput `json:"translations"`
	ProductTypeId  int                `json:"product_type_id"`
	ManufacturerId int                `json:"manufacturer_id"`
}

type PatchInput struct {
	Name             *string
	Description      *string
	ShortDescription *string
	ProductTypeId    *int
	ManufacturerId   *int
}
