package categories_service

import (
	"context"
	"github.com/google/uuid"
	category_model "kadabra/internal/features/categories/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *category_model.Category) (*category_model.Category, error)
	GetAll(ctx context.Context) ([]*category_model.Category, error)
	GetById(ctx context.Context, id uuid.UUID) (*category_model.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, category *category_model.CategoryPatch) (*category_model.Category, error)
}
