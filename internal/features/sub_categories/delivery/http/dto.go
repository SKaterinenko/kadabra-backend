package sub_categories_http

import sub_categories_service "kadabra/internal/features/sub_categories/service"

type TranslationInput struct {
	LanguageCode string
	Name         string
}

type Translations []TranslationInput

type createDTO struct {
	sub_categories_service.CreateInput
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}
