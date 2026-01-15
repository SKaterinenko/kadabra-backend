package manufacturers_service

import (
	"context"
	manufacturers_model "kadabra/internal/features/manufacturers/model"
)

type ManufacturerRepository interface {
	Create(ctx context.Context, manufacturer *CreateInput) (*manufacturers_model.ManufacturerWithTranslations, error)
	GetAll(ctx context.Context, lang string) ([]*manufacturers_model.Manufacturer, error)
	GetById(ctx context.Context, id int, lang string) (*manufacturers_model.Manufacturer, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, manufacturer *manufacturers_model.ManufacturerPatch) (*manufacturers_model.Manufacturer, error)
}
