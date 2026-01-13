package products_type_service

import (
	"context"

	"kadabra/internal/features/products_type/model"
	sub_categories_model "kadabra/internal/features/sub_categories/model"
)

type ProductsTypeRepository interface {
	Create(ctx context.Context, productsType *products_type_model.ProductsType) (*products_type_model.ProductsType, error)
	GetAll(ctx context.Context) ([]*products_type_model.ProductsType, error)
	GetById(ctx context.Context, id int) (*products_type_model.ProductsType, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, productsType *products_type_model.ProductsTypePatch) (*products_type_model.ProductsType, error)
	GetProductsTypeByCategorySlug(ctx context.Context, slug string) ([]*sub_categories_model.SubCategoryWithProductsType, error)
}
