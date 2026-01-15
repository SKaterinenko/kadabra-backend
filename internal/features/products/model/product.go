package products_model

import (
	"time"
)

type ProductTranslate struct {
	Id               int       `json:"id" db:"id"`
	ProductId        int       `json:"product_id" db:"product_id"`
	LanguageCode     string    `json:"language_code" db:"language_code"`
	Name             string    `json:"name" db:"name"`
	Slug             string    `json:"slug" db:"slug"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Description      string    `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type ProductWithoutTranslations struct {
	Id             int       `json:"id" db:"id"`
	ProductsTypeId int       `json:"products_type_id" db:"products_type_id"`
	ManufacturerId int       `json:"manufacturer_id" db:"manufacturer_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type ProductWithTranslations struct {
	Id             int                 `json:"id" db:"id"`
	ProductsTypeId int                 `json:"products_type_id" db:"products_type_id"`
	ManufacturerId int                 `json:"manufacturer_id" db:"manufacturer_id"`
	Translations   []*ProductTranslate `json:"translations" db:"translations"`
	CreatedAt      time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" db:"updated_at"`
}

type Product struct {
	Id               int       `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Slug             string    `json:"slug" db:"slug"`
	ProductsTypeId   int       `json:"products_type_id" db:"products_type_id"`
	ManufacturerId   int       `json:"manufacturer_id" db:"manufacturer_id"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Description      string    `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type ProductPatch struct {
	Name             *string `json:"name,omitempty"`
	Description      *string `json:"description,omitempty"`
	ShortDescription *string `json:"short_description,omitempty"`
	ProductsTypeId   *int    `json:"products_type_id,omitempty"`
	ManufacturerId   *int    `json:"manufacturer_id,omitempty"`
}

func NewProductPatch(name, description, shortDescription string, productsTypeId, manufacturerId int) *ProductPatch {
	return &ProductPatch{
		Name:             &name,
		ShortDescription: &shortDescription,
		Description:      &description,
		ProductsTypeId:   &productsTypeId,
		ManufacturerId:   &manufacturerId,
	}
}
