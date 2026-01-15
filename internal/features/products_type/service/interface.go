package products_type_service

import (
	"context"

	products_type_model "kadabra/internal/features/products_type/model"
	sub_categories_model "kadabra/internal/features/sub_categories/model"
)

type ProductsTypeRepository interface {
	Create(ctx context.Context, productsType *CreateInput) (*products_type_model.ProductsTypeWithTranslations, error)
	GetAll(ctx context.Context, lang string) ([]*products_type_model.ProductsType, error)
	GetById(ctx context.Context, id int, lang string) (*products_type_model.ProductsType, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, productsType *products_type_model.ProductsTypePatch) (*products_type_model.ProductsType, error)
	GetProductsTypeByCategorySlug(ctx context.Context, slug string) ([]*sub_categories_model.SubCategoryWithProductsType, error)
}
