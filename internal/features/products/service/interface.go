package products_service

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/features/products/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *products_model.Product) (*products_model.Product, error)
	GetAll(ctx context.Context) ([]*products_model.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*products_model.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, product *products_model.ProductPatch) (*products_model.Product, error)
	GetByCategoryIds(ctx context.Context, categoryIds []uuid.UUID) ([]*products_model.Product, error)
}
