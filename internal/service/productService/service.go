package productService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type Service struct {
	repo ProductRepository
}

func NewService(repo ProductRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, category *CreateInput) (*model.Product, error) {
	newProduct := model.NewProduct(
		category.Name,
		category.Description,
		category.ShortDescription,
		category.ProductsTypeId,
		category.ManufacturerId)

	out, err := s.repo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*model.Product, error) {
	out, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*model.Product, error) {
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

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*model.Product, error) {
	newPatch := model.NewProductPatch(
		*update.Name,
		*update.Description,
		*update.ShortDescription,
		*update.ProductsTypeId,
		*update.ManufacturerId)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetByCategoryIds(ctx context.Context, categoryIds []uuid.UUID) ([]*model.Product, error) {
	out, err := s.repo.GetByCategoryIds(ctx, categoryIds)
	if err != nil {
		return nil, err
	}
	return out, nil
}
