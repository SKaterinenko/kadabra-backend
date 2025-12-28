package model

import (
	"github.com/google/uuid"
	"time"
)

type SubCategory struct {
	Id         uuid.UUID `json:"id"`
	CategoryId uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewSubCategory(name string, categoryId uuid.UUID) *SubCategory {
	return &SubCategory{
		Id:         uuid.New(),
		CategoryId: categoryId,
		Name:       name,
	}
}

type SubCategoryPatch struct {
	Name *string `json:"name,omitempty"`
}

func NewSubCategoryPatch(name string) *SubCategoryPatch {
	return &SubCategoryPatch{
		Name: &name,
	}
}

type SubCategoryWithProductsType struct {
	SubCategory
	ProductsType []*ProductsType `json:"products_type"`
}
