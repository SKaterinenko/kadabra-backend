package manufacturerService

import (
	"context"
	"github.com/google/uuid"
	"kadabra/internal/model"
)

type ManufacturerRepository interface {
	Create(ctx context.Context, manufacturer *model.Manufacturer) error
	GetAll(ctx context.Context) ([]*model.Manufacturer, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.Manufacturer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, manufacturer *model.ManufacturerPatch) (*model.Manufacturer, error)
}
