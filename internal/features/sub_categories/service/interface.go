package sub_categories_service

import (
	"context"

	"kadabra/internal/features/sub_categories/model"
)

type SubCategoryRepository interface {
	Create(ctx context.Context, req *CreateInput) (*sub_categories_model.SubCategoryWithTranslations, error)
	GetAll(ctx context.Context, lang string) ([]*sub_categories_model.SubCategory, error)
	GetById(ctx context.Context, id int, lang string) (*sub_categories_model.SubCategory, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, subCategory *sub_categories_model.SubCategoryPatch) (*sub_categories_model.SubCategoryWithTranslations, error)
}
