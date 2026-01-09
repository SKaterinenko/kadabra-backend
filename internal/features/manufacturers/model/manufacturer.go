package manufacturers_model

import (
	"github.com/google/uuid"
	"time"
)

type Manufacturer struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewManufacturer(name string) *Manufacturer {
	return &Manufacturer{
		Id:   uuid.New(),
		Name: name,
	}
}

type ManufacturerPatch struct {
	Name *string `json:"name,omitempty"`
}

func NewManufacturerPatch(name string) *ManufacturerPatch {
	return &ManufacturerPatch{
		Name: &name,
	}
}
