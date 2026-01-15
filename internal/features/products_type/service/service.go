package products_type_service

import (
	"context"
	products_type_model "kadabra/internal/features/products_type/model"
	sub_categories_model "kadabra/internal/features/sub_categories/model"
)

type Service struct {
	repo ProductsTypeRepository
}

func NewService(repo ProductsTypeRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, productsType *CreateInput) (*products_type_model.ProductsTypeWithTranslations, error) {
	out, err := s.repo.Create(ctx, productsType)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context, lang string) ([]*products_type_model.ProductsType, error) {
	out, err := s.repo.GetAll(ctx, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id int, lang string) (*products_type_model.ProductsType, error) {
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

func (s *Service) Patch(ctx context.Context, id int, update *PatchInput) (*products_type_model.ProductsType, error) {
	newPatch := products_type_model.NewProductsTypePatch(*update.Name)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetProductsTypeByCategorySlug(ctx context.Context, slug string) ([]*sub_categories_model.SubCategoryWithProductsType, error) {
	out, err := s.repo.GetProductsTypeByCategorySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return out, nil
}
