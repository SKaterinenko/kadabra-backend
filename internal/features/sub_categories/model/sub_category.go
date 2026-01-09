package sub_categories_model

import (
	"github.com/google/uuid"
	products_type_model "kadabra/internal/features/products_type/model"
	"time"
)

type SubCategory struct {
	Id         uuid.UUID `json:"id" db:"id"`
	CategoryId uuid.UUID `json:"category_id" db:"category_id"`
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func NewSubCategory(name string, categoryId uuid.UUID) *SubCategory {
	return &SubCategory{
		Id:         uuid.New(),
		CategoryId: categoryId,
		Name:       name,
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
