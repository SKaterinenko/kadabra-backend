package subCategoryService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type SubCategoryRepository interface {
	Create(ctx context.Context, subCategory *model.SubCategory) (*model.SubCategory, error)
	GetAll(ctx context.Context) ([]*model.SubCategory, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.SubCategory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, subCategory *model.SubCategoryPatch) (*model.SubCategory, error)
}
