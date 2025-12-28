package categoryService

type CreateInput struct {
	Name string
}

type PatchInput struct {
	Name *string
}
