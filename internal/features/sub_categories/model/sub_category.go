package sub_categories_model

import (
	products_type_model "kadabra/internal/features/products_type/model"
	"time"
)

type TranslationInput struct {
	LanguageCode string
	Name         string
}

type Translations []TranslationInput

type SubCategoryWithoutTranslations struct {
	Id         int       `json:"id" db:"id"`
	CategoryId int       `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type SubCategoryWithTranslations struct {
	Id           int                     `json:"id" db:"id"`
	CategoryId   int                     `json:"category_id" db:"category_id"`
	Translations []*SubCategoryTranslate `json:"translations" db:"translations"`
	CreatedAt    time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at" db:"updated_at"`
}

type SubCategory struct {
	Id         int       `json:"id" db:"id"`
	CategoryId int       `json:"category_id" db:"category_id"`
	Name       string    `json:"name" db:"name"`
	Slug       string    `json:"slug" db:"slug"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type SubCategoryTranslate struct {
	Id            int       `json:"id" db:"id"`
	SubCategoryId int       `json:"sub_category_id" db:"sub_category_id"`
	Name          string    `json:"name" db:"name"`
	Slug          string    `json:"slug" db:"slug"`
	LanguageCode  string    `json:"language_code" db:"language_code"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type SubCategoryPatch struct {
	Name *string `json:"name,omitempty" db:"name"`
}

func NewSubCategoryPatch(name string) *SubCategoryPatch {
	return &SubCategoryPatch{
		Name: &name,
	}
}

type SubCategoryWithProductsType struct {
	Id           int                                 `json:"id" db:"id"`
	CategoryId   int                                 `json:"category_id" db:"category_id"`
	Name         string                              `json:"name" db:"name"`
	Slug         string                              `json:"slug" db:"slug"`
	ProductsType []*products_type_model.ProductsType `json:"products_type" db:"products_type"`
	CreatedAt    time.Time                           `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time                           `json:"updated_at" db:"updated_at"`
}
