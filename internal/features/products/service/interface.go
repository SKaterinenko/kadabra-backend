package products_service

import (
	"context"

	"kadabra/internal/features/products/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *products_model.Product) (*products_model.Product, error)
	GetAll(ctx context.Context) ([]*products_model.Product, error)
	GetById(ctx context.Context, id int) (*products_model.Product, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, product *products_model.ProductPatch) (*products_model.Product, error)
	GetByCategoryIds(ctx context.Context, categoryIds []int) ([]*products_model.Product, error)
}
