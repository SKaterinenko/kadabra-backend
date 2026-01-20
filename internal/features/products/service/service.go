package products_service

import (
	"context"
	products_model "kadabra/internal/features/products/model"
)

type Service struct {
	repo ProductRepository
}

func NewService(repo ProductRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, product *CreateInput) (*products_model.ProductWithTranslations, error) {
	out, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context, lang string) ([]*products_model.Product, error) {
	out, err := s.repo.GetAll(ctx, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id int, lang string) (*products_model.Product, error) {
	out, err := s.repo.GetById(ctx, id, lang)
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
		*update.ProductTypeId,
		*update.ManufacturerId)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetByCategoryIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error) {
	if len(categoryIds) == 0 {
		return []*products_model.Product{}, nil
	}

	out, err := s.repo.GetByCategoryIds(ctx, categoryIds, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetByProductsTypeIds(ctx context.Context, productsTypeIds []int, lang string) ([]*products_model.Product, error) {
	if len(productsTypeIds) == 0 {
		return []*products_model.Product{}, nil
	}

	out, err := s.repo.GetByProductsTypeIds(ctx, productsTypeIds, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetByCategorySlug(ctx context.Context, lang, slug string) ([]*products_model.Product, error) {
	if slug == "" {
		return []*products_model.Product{}, nil
	}
	products, err := s.repo.GetByCategorySlug(ctx, lang, slug)
	if err != nil {
		return nil, err
	}

	return products, nil
}
