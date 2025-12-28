package productsTypeService

import (
	"github.com/google/uuid"
)

type CreateInput struct {
	Name          string
	SubCategoryId uuid.UUID
}

type PatchInput struct {
	Name *string
}
