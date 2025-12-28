package categoryService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type Service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, category *CreateInput) (*model.Category, error) {
	newCategory := model.NewCategory(category.Name)

	err := s.repo.Create(ctx, newCategory)
	if err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
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

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*model.Category, error) {
	newPatch := model.NewCategoryPatch(*update.Name)
	patch, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return patch, nil
}
