package manufacturers_model

import (
	"time"
)

type ManufacturerTranslate struct {
	Id             int       `json:"id" db:"id"`
	ManufacturerId int       `json:"manufacturer_id" db:"manufacturer_id"`
	LanguageCode   string    `json:"language_code" db:"language_code"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type ManufacturerWithoutTranslations struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ManufacturerWithTranslations struct {
	Id           int                      `json:"id" db:"id"`
	Name         string                   `json:"name" db:"name"`
	Slug         string                   `json:"slug" db:"slug"`
	Translations []*ManufacturerTranslate `json:"translations" db:"translations"`
	CreatedAt    time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at" db:"updated_at"`
}

type Manufacturer struct {
	Id          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ManufacturerPatch struct {
	Name *string `json:"name,omitempty"`
}

func NewManufacturerPatch(name string) *ManufacturerPatch {
	return &ManufacturerPatch{
		Name: &name,
	}
}
