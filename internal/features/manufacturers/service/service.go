package manufacturers_service

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/features/manufacturers/model"
)

type Service struct {
	repo ManufacturerRepository
}

func NewService(repo ManufacturerRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, category *CreateInput) (*manufacturers_model.Manufacturer, error) {
	newManufacturer := manufacturers_model.NewManufacturer(category.Name)

	out, err := s.repo.Create(ctx, newManufacturer)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*manufacturers_model.Manufacturer, error) {
	out, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*manufacturers_model.Manufacturer, error) {
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

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*manufacturers_model.Manufacturer, error) {
	newPatch := manufacturers_model.NewManufacturerPatch(*update.Name)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}
