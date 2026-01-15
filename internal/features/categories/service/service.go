package categories_service

import (
	"context"
	"fmt"
	categories_model "kadabra/internal/features/categories/model"
)

type Service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateInput) (*categories_model.CategoryWithTranslations, error) {
	// Минимум один перевод
	if len(req.Translations) == 0 {
		return nil, fmt.Errorf("at least one translation required")
	}

	return s.repo.Create(ctx, req)
}

func (s *Service) GetAll(ctx context.Context, language string) ([]*categories_model.Category, error) {
	return s.repo.GetAll(ctx, language)
}

func (s *Service) GetById(ctx context.Context, id int, language string) (*categories_model.Category, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid category id")
	}

	return s.repo.GetById(ctx, id, language)
}

func (s *Service) GetBySlug(ctx context.Context, slug, language string) (*categories_model.Category, error) {
	if slug == "" {
		return nil, fmt.Errorf("slug cannot be empty")
	}

	return s.repo.GetBySlug(ctx, slug, language)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

//func (s *Service) Patch(ctx context.Context, id int, update *PatchInput) (*category_model.Category, error) {
//	newPatch := category_model.NewCategoryPatch(*update.Name)
//	out, err := s.repo.Patch(ctx, id, newPatch)
//	if err != nil {
//		return nil, err
//	}
//	return out, nil
//}
