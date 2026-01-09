package sub_categories_service

import (
	"github.com/google/uuid"
)

type CreateInput struct {
	Name       string
	CategoryId uuid.UUID
}

type PatchInput struct {
	Name *string
}
