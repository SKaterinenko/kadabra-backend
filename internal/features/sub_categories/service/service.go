package sub_categories_service

import (
	"context"

	"github.com/gosimple/slug"
	"kadabra/internal/features/sub_categories/model"
)

type Service struct {
	repo SubCategoryRepository
}

func NewService(repo SubCategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, subCategory *CreateInput) (*sub_categories_model.SubCategory, error) {
	slugText := slug.Make(subCategory.Name)
	newSubCategory := sub_categories_model.NewSubCategory(subCategory.Name, slugText, subCategory.CategoryId)

	out, err := s.repo.Create(ctx, newSubCategory)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*sub_categories_model.SubCategory, error) {
	out, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id int) (*sub_categories_model.SubCategory, error) {
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

func (s *Service) Patch(ctx context.Context, id int, update *PatchInput) (*sub_categories_model.SubCategory, error) {
	newPatch := sub_categories_model.NewSubCategoryPatch(*update.Name)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}
