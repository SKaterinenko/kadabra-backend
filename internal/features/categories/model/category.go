package categories_model

import (
	"time"
)

type CategoryWithoutTranslations struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CategoryTranslate struct {
	Id           int       `db:"id" json:"id"`
	CategoryID   int       `db:"category_id" json:"category_id"`
	LanguageCode string    `db:"language_code" json:"language_code"`
	Name         string    `db:"name" json:"name"`
	Slug         string    `db:"slug" json:"slug"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CategoryWithTranslations struct {
	Id           int                  `json:"id"`
	Translations []*CategoryTranslate `json:"translations"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type CategoryPatch struct {
	Name *string `json:"name,omitempty" db:"name"`
}
