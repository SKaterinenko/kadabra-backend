package sub_categories_http

import (
	"context"
	sub_categories_model "kadabra/internal/features/sub_categories/model"
	sub_categories_service "kadabra/internal/features/sub_categories/service"
)

type SubCategoriesService interface {
	Create(ctx context.Context, subCategory *sub_categories_service.CreateInput) (*sub_categories_model.SubCategoryWithTranslations, error)
	GetAll(ctx context.Context, lang string) ([]*sub_categories_model.SubCategory, error)
	GetById(ctx context.Context, id int, lang string) (*sub_categories_model.SubCategory, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, update *sub_categories_service.PatchInput) (*sub_categories_model.SubCategoryWithTranslations, error)
}
