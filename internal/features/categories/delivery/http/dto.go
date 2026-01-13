package categories_http

import categories_service "kadabra/internal/features/categories/service"

type TranslationInput struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

type createDTO struct {
	categories_service.CreateInput
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}
