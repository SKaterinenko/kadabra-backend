package productsTypeService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type Service struct {
	repo ProductsTypeRepository
}

func NewService(repo ProductsTypeRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, subCategory *CreateInput) (*model.ProductsType, error) {
	newSubCategory := model.NewProductsType(subCategory.Name, subCategory.SubCategoryId)

	out, err := s.repo.Create(ctx, newSubCategory)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*model.ProductsType, error) {
	out, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*model.ProductsType, error) {
	out, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*model.ProductsType, error) {
	newPatch := model.NewProductsTypePatch(*update.Name)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetProductsTypeByCategoryId(ctx context.Context, id uuid.UUID) ([]*model.SubCategoryWithProductsType, error) {
	out, err := s.repo.GetProductsTypeByCategoryId(ctx, id)
	if err != nil {
		return nil, err
	}
	return out, nil
}
