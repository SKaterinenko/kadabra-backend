package manufacturers_model

import (
	"time"
)

type Manufacturer struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewManufacturer(name, slug string) *Manufacturer {
	return &Manufacturer{
		Name: name,
		Slug: slug,
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
