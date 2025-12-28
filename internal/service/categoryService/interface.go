package categoryService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetAll(ctx context.Context) ([]*model.Category, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, category *model.CategoryPatch) (*model.Category, error)
}
