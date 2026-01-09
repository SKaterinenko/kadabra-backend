package sub_categories_service

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/features/sub_categories/model"
)

type SubCategoryRepository interface {
	Create(ctx context.Context, subCategory *sub_categories_model.SubCategory) (*sub_categories_model.SubCategory, error)
	GetAll(ctx context.Context) ([]*sub_categories_model.SubCategory, error)
	GetById(ctx context.Context, id uuid.UUID) (*sub_categories_model.SubCategory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, subCategory *sub_categories_model.SubCategoryPatch) (*sub_categories_model.SubCategory, error)
}
