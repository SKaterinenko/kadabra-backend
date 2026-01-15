package categories_http

import (
	"context"
	categories_model "kadabra/internal/features/categories/model"
	categories_service "kadabra/internal/features/categories/service"
)

type CategoryService interface {
	Create(ctx context.Context, req *categories_service.CreateInput) (*categories_model.CategoryWithTranslations, error)
	GetAll(ctx context.Context, language string) ([]*categories_model.Category, error)
	GetById(ctx context.Context, id int, language string) (*categories_model.Category, error)
	GetBySlug(ctx context.Context, slug, language string) (*categories_model.Category, error)
	Delete(ctx context.Context, id int) error
	//Patch(ctx context.Context, id int, update *categories_service.PatchInput) (*categories_model.Category, error)
}
