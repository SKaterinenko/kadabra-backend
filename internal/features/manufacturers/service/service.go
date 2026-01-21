package manufacturers_service

import (
	"context"
	manufacturers_model "kadabra/internal/features/manufacturers/model"
)

type Service struct {
	repo ManufacturerRepository
}

func NewService(repo ManufacturerRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, manufacturer *CreateInput) (*manufacturers_model.ManufacturerWithTranslations, error) {
	out, err := s.repo.Create(ctx, manufacturer)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context, lang string) ([]*manufacturers_model.Manufacturer, error) {
	out, err := s.repo.GetAll(ctx, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id int, lang string) (*manufacturers_model.Manufacturer, error) {
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

func (s *Service) Patch(ctx context.Context, id int, update *PatchInput) (*manufacturers_model.ManufacturerWithTranslations, error) {
	out, err := s.repo.Patch(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return out, nil
}
