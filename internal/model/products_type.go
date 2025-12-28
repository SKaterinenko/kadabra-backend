package model

import (
	"github.com/google/uuid"
	"time"
)

type ProductsType struct {
	Id            uuid.UUID `json:"id"`
	SubCategoryId uuid.UUID `json:"sub_category_id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewProductsType(name string, subCategoryId uuid.UUID) *ProductsType {
	return &ProductsType{
		Id:            uuid.New(),
		SubCategoryId: subCategoryId,
		Name:          name,
	}
}

type ProductsTypePatch struct {
	Name *string `json:"name,omitempty"`
}

func NewProductsTypePatch(name string) *SubCategoryPatch {
	return &SubCategoryPatch{
		Name: &name,
	}
}
