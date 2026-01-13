package products_type_http

import (
	"context"

	products_type_model "kadabra/internal/features/products_type/model"
	products_type_service "kadabra/internal/features/products_type/service"
	sub_categories_model "kadabra/internal/features/sub_categories/model"
)

type ProductsTypeInterface interface {
	Create(ctx context.Context, subCategory *products_type_service.CreateInput) (*products_type_model.ProductsType, error)
	GetAll(ctx context.Context) ([]*products_type_model.ProductsType, error)
	GetById(ctx context.Context, id int) (*products_type_model.ProductsType, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, update *products_type_service.PatchInput) (*products_type_model.ProductsType, error)
	GetProductsTypeByCategorySlug(ctx context.Context, slug string) ([]*sub_categories_model.SubCategoryWithProductsType, error)
}
