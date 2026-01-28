package products_service

import (
	"context"

	products_model "kadabra/internal/features/products/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *CreateInput) (*products_model.ProductWithTranslations, error)
	GetAll(
		ctx context.Context,
		lang string,
		categories, types, manufacturers []int,
		limit, offset int,
	) (*products_model.Products, error)
	GetById(ctx context.Context, id int, lang string) (*products_model.Product, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, product *products_model.ProductPatch) (*products_model.Product, error)
	GetByCategoryIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error)
	GetByProductsTypeIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error)
	GetByCategorySlug(ctx context.Context, lang, slug string) ([]*products_model.Product, error)
	GetByManufacturersIds(ctx context.Context, ids []int, lang string) ([]*products_model.Product, error)
	GetBySlug(ctx context.Context, slug, lang string) (*products_model.ProductWithParents, error)
	CreateProductVariations(ctx context.Context, req *VariationReq) (*products_model.ProductVariation, error)
	DeleteProductVariation(ctx context.Context, id int) error
}
