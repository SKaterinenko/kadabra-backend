package productsTypeService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type ProductsTypeRepository interface {
	Create(ctx context.Context, productsType *model.ProductsType) (*model.ProductsType, error)
	GetAll(ctx context.Context) ([]*model.ProductsType, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.ProductsType, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, productsType *model.SubCategoryPatch) (*model.ProductsType, error)
	GetProductsTypeByCategoryId(ctx context.Context, id uuid.UUID) ([]*model.SubCategoryWithProductsType, error)
}
