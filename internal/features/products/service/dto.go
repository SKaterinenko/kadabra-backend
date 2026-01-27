package products_service

import "mime/multipart"

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

type VariationInput struct {
	ProductId int                   `json:"product_id"`
	Image     *multipart.FileHeader `json:"image"`
	Price     int                   `json:"price"`
}

type VariationReq struct {
	ProductId int    `json:"product_id" db:"product_id"`
	Image     string `json:"image" db:"image"`
	Price     int    `json:"price" db:"price"`
}
