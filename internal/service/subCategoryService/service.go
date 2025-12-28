package subCategoryService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type Service struct {
	repo SubCategoryRepository
}

func NewService(repo SubCategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, category *CreateInput) (*model.SubCategory, error) {
	newCategory := model.NewSubCategory(category.Name, category.CategoryId)

	err := s.repo.Create(ctx, newCategory)
	if err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*model.SubCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*model.SubCategory, error) {
	byId, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return byId, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*model.SubCategory, error) {
	newPatch := model.NewSubCategoryPatch(*update.Name)
	patch, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return patch, nil
}
