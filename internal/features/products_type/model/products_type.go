package products_type_model

import (
	"time"
)

type ProductsType struct {
	Id            int       `json:"id" db:"id"`
	SubCategoryId int       `json:"sub_category_id" db:"sub_category_id"`
	Name          string    `json:"name" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func NewProductsType(name string, subCategoryId int) *ProductsType {
	return &ProductsType{
		SubCategoryId: subCategoryId,
		Name:          name,
	}
}

type ProductsTypePatch struct {
	Name *string `json:"name,omitempty" db:"name"`
}

func NewProductsTypePatch(name string) *ProductsTypePatch {
	return &ProductsTypePatch{
		Name: &name,
	}
}
