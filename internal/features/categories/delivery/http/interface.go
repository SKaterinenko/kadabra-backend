package categories_http

import (
	"context"
	"github.com/google/uuid"
	category_model "kadabra/internal/features/categories/model"
	category_service "kadabra/internal/features/categories/service"
)

type CategoryService interface {
	Create(ctx context.Context, category *category_service.CreateInput) (*category_model.Category, error)
	GetAll(ctx context.Context) ([]*category_model.Category, error)
	GetById(ctx context.Context, id uuid.UUID) (*category_model.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, update *category_service.PatchInput) (*category_model.Category, error)
}
