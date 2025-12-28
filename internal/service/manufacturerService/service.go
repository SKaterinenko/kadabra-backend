package manufacturerService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type Service struct {
	repo ManufacturerRepository
}

func NewService(repo ManufacturerRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, category *CreateInput) (*model.Manufacturer, error) {
	newManufacturer := model.NewManufacturer(category.Name)

	err := s.repo.Create(ctx, newManufacturer)
	if err != nil {
		return nil, err
	}

	return newManufacturer, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*model.Manufacturer, error) {
	manufacturers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return manufacturers, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*model.Manufacturer, error) {
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

func (s *Service) Patch(ctx context.Context, id uuid.UUID, update *PatchInput) (*model.Manufacturer, error) {
	newPatch := model.NewManufacturerPatch(*update.Name)
	patch, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return patch, nil
}
