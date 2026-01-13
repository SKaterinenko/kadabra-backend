package categories_model

import (
	"time"
)

type Category struct {
	Id        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CategoryTranslation struct {
	Id           int       `db:"id" json:"id"`
	CategoryID   int       `db:"category_id" json:"category_id"`
	LanguageCode string    `db:"language_code" json:"language_code"`
	Name         string    `db:"name" json:"name"`
	Slug         string    `db:"slug" json:"slug"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCategoryRequest struct {
	Translations map[string]TranslationInput `json:"translations" binding:"required"`
}

type CategoryResponse struct {
	Id        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TranslationInput - данные для перевода
type TranslationInput struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// CategoryWithTranslations - категория со всеми переводами
type CategoryWithTranslations struct {
	Id           int                            `json:"id"`
	Translations map[string]CategoryTranslation `json:"translations"`
	CreatedAt    time.Time                      `json:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at"`
}

type CategoryPatch struct {
	Name *string `json:"name,omitempty" db:"name"`
}

func NewCategoryPatch(name string) *CategoryPatch {
	return &CategoryPatch{
		Name: &name,
	}
}
