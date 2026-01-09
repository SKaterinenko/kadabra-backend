package manufacturers_service

type CreateInput struct {
	Name string
}

type PatchInput struct {
	Name *string
}
