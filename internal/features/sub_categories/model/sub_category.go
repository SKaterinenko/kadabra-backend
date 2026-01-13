package sub_categories_model

import (
	products_type_model "kadabra/internal/features/products_type/model"
	"time"
)

type SubCategory struct {
	Id         int       `json:"id" db:"id"`
	CategoryId int       `json:"category_id" db:"category_id"`
	Name       string    `json:"name" db:"name"`
	Slug       string    `json:"slug" db:"slug"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func NewSubCategory(name, slug string, categoryId int) *SubCategory {
	return &SubCategory{
		CategoryId: categoryId,
		Name:       name,
		Slug:       slug,
	}
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
	SubCategory
	ProductsType []*products_type_model.ProductsType `json:"products_type" db:"products_type"`
}
