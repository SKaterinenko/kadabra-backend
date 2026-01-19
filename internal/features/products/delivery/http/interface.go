package products_http

import (
	"context"
	products_model "kadabra/internal/features/products/model"
	products_service "kadabra/internal/features/products/service"
)

type ProductsService interface {
	Create(ctx context.Context, product *products_service.CreateInput) (*products_model.ProductWithTranslations, error)
	GetAll(ctx context.Context, lang string) ([]*products_model.Product, error)
	GetById(ctx context.Context, id int, lang string) (*products_model.Product, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, update *products_service.PatchInput) (*products_model.Product, error)
	GetByCategoryIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error)
	GetByProductsTypeIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error)
}
