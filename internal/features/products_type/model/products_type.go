package products_type_model

import (
	"time"
)

type ProductsTypeTranslate struct {
	Id            int       `json:"id" db:"id"`
	ProductTypeId int       `json:"product_type_id" db:"product_type_id"`
	LanguageCode  string    `json:"language_code" db:"language_code"`
	Name          string    `json:"name" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type ProductsTypeWithoutTranslations struct {
	Id            int       `json:"id" db:"id"`
	SubCategoryId int       `json:"sub_category_id" db:"sub_category_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type ProductsTypeWithTranslations struct {
	Id            int                      `json:"id" db:"id"`
	SubCategoryId int                      `json:"sub_category_id" db:"sub_category_id"`
	Translations  []*ProductsTypeTranslate `json:"translations" db:"translations"`
	CreatedAt     time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at" db:"updated_at"`
}

type ProductsType struct {
	Id            int       `json:"id" db:"id"`
	SubCategoryId int       `json:"sub_category_id" db:"sub_category_id"`
	Name          string    `json:"name" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type ProductsTypePatch struct {
	Name *string `json:"name,omitempty" db:"name"`
}

func NewProductsTypePatch(name string) *ProductsTypePatch {
	return &ProductsTypePatch{
		Name: &name,
	}
}
