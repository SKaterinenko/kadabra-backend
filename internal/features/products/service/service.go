package products_service

import (
	"context"

	"github.com/gosimple/slug"
	"kadabra/internal/features/products/model"
)

type Service struct {
	repo ProductRepository
}

func NewService(repo ProductRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, product *CreateInput) (*products_model.Product, error) {
	slugText := slug.Make(product.Name)
	newProduct := products_model.NewProduct(
		product.Name,
		slugText,
		product.Description,
		product.ShortDescription,
		product.ProductsTypeId,
		product.ManufacturerId)

	out, err := s.repo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*products_model.Product, error) {
	out, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id int) (*products_model.Product, error) {
	out, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Patch(ctx context.Context, id int, update *PatchInput) (*products_model.Product, error) {
	newPatch := products_model.NewProductPatch(
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

func (s *Service) GetByCategoryIds(ctx context.Context, categoryIds []int) ([]*products_model.Product, error) {
	out, err := s.repo.GetByCategoryIds(ctx, categoryIds)
	if err != nil {
		return nil, err
	}
	return out, nil
}
