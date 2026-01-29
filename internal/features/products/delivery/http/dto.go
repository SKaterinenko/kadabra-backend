package products_http

import (
	"github.com/shopspring/decimal"
	"mime/multipart"
)

type TranslationInput struct {
	LanguageCode     string `json:"language_code" validate:"required"`
	Name             string `json:"name" validate:"required"`
	ShortDescription string `json:"short_description" validate:"required"`
	Description      string `json:"description" validate:"required"`
}

type createDTO struct {
	Translations   []TranslationInput `json:"translations" validate:"required,min=1,dive"`
	ProductTypeId  int                `json:"product_type_id" validate:"required"`
	ManufacturerId int                `json:"manufacturer_id" validate:"required"`
}

type patchDTO struct {
	Name             *string `json:"name,omitempty" validate:"omitempty"`
	Description      *string `json:"description,omitempty" validate:"omitempty"`
	ShortDescription *string `json:"short_description,omitempty" validate:"omitempty"`
	ProductTypeId    *int    `json:"product_type_id,omitempty" validate:"omitempty"`
	ManufacturerId   *int    `json:"manufacturer_id,omitempty" validate:"omitempty"`
}

type getByIdsDTO struct {
	Ids []int `json:"ids" validate:"required"`
}

type createVariationDTO struct {
	ProductId int                   `form:"product_id" validate:"required,gt=0"`
	Image     *multipart.FileHeader `form:"image" validate:"required,image_type"`
	Price     decimal.Decimal       `form:"price" validate:"required,gte=0"`
}
