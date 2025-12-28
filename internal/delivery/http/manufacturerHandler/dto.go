package manufacturerHandler

type createDTO struct {
	Name string `json:"name" validate:"required"`
}

type patchDTO struct {
	Name *string `json:"name,omitempty" validate:"omitempty"`
}
