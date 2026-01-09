package products_model

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	ProductsTypeId   uuid.UUID `json:"products_type_id" db:"products_type_id"`
	ManufacturerId   uuid.UUID `json:"manufacturer_id" db:"manufacturer_id"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Description      string    `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func NewProduct(name, description, shortDescription string, productsTypeId, manufacturerId uuid.UUID) *Product {
	return &Product{
		Id:               uuid.New(),
		Name:             name,
		ProductsTypeId:   productsTypeId,
		ManufacturerId:   manufacturerId,
		ShortDescription: shortDescription,
		Description:      description,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

type ProductPatch struct {
	Name             *string    `json:"name,omitempty"`
	Description      *string    `json:"description,omitempty"`
	ShortDescription *string    `json:"short_description,omitempty"`
	ProductsTypeId   *uuid.UUID `json:"products_type_id,omitempty"`
	ManufacturerId   *uuid.UUID `json:"manufacturer_id,omitempty"`
}

func NewProductPatch(name, description, shortDescription string, productsTypeId, manufacturerId uuid.UUID) *ProductPatch {
	return &ProductPatch{
		Name:             &name,
		ShortDescription: &shortDescription,
		Description:      &description,
		ProductsTypeId:   &productsTypeId,
		ManufacturerId:   &manufacturerId,
	}
}
