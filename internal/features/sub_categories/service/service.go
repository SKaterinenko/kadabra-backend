package sub_categories_service

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/features/sub_categories/model"
)

type Service struct {
	repo SubCategoryRepository
}

func NewService(repo SubCategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, category *CreateInput) (*sub_categories_model.SubCategory, error) {
	newCategory := sub_categories_model.NewSubCategory(category.Name, category.CategoryId)

	out, err := s.repo.Create(ctx, newCategory)
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

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*sub_categories_model.SubCategory, error) {
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

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*sub_categories_model.SubCategory, error) {
	newPatch := sub_categories_model.NewSubCategoryPatch(*update.Name)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}
