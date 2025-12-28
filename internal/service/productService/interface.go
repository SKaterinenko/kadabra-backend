package productService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetAll(ctx context.Context) ([]*model.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, product *model.ProductPatch) (*model.Product, error)
	GetByCategoryIds(ctx context.Context, categoryIds []uuid.UUID) ([]*model.Product, error)
}
