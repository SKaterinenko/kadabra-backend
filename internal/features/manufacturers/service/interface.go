package manufacturers_service

import (
	"context"
	"kadabra/internal/features/manufacturers/model"
)

type ManufacturerRepository interface {
	Create(ctx context.Context, manufacturer *manufacturers_model.Manufacturer) (*manufacturers_model.Manufacturer, error)
	GetAll(ctx context.Context) ([]*manufacturers_model.Manufacturer, error)
	GetById(ctx context.Context, id int) (*manufacturers_model.Manufacturer, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, manufacturer *manufacturers_model.ManufacturerPatch) (*manufacturers_model.Manufacturer, error)
}
