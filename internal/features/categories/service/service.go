package categories_service

import (
	"context"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	category_model "kadabra/internal/features/categories/model"
)

type Service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, category *CreateInput) (*category_model.Category, error) {
	slugText := slug.Make(category.Name)
	newCategory := category_model.NewCategory(category.Name, slugText)
	out, err := s.repo.Create(ctx, newCategory)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*category_model.Category, error) {
	out, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*category_model.Category, error) {
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

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*category_model.Category, error) {
	newPatch := category_model.NewCategoryPatch(*update.Name)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}
