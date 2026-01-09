package products_type_service

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/features/products_type/model"
	sub_categories_model "kadabra/internal/features/sub_categories/model"
)

type ProductsTypeRepository interface {
	Create(ctx context.Context, productsType *products_type_model.ProductsType) (*products_type_model.ProductsType, error)
	GetAll(ctx context.Context) ([]*products_type_model.ProductsType, error)
	GetById(ctx context.Context, id uuid.UUID) (*products_type_model.ProductsType, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, productsType *products_type_model.ProductsTypePatch) (*products_type_model.ProductsType, error)
	GetProductsTypeByCategoryId(ctx context.Context, id uuid.UUID) ([]*sub_categories_model.SubCategoryWithProductsType, error)
}
