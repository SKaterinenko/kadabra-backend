package categories_service

import (
	"context"
	categories_model "kadabra/internal/features/categories/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, req *CreateInput) (*categories_model.CategoryWithTranslations, error)
	GetAll(ctx context.Context, language string) ([]*categories_model.Category, error)
	GetById(ctx context.Context, id int, language string) (*categories_model.Category, error)
	GetBySlug(ctx context.Context, slug, language string) (*categories_model.Category, error)
	Delete(ctx context.Context, id int) error
	//Patch(ctx context.Context, id int, category *categories_model.CategoryPatch) (*categories_model.Category, error)
}
