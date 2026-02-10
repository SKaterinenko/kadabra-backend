package reviews_http

type createDTO struct {
	ProductId   int64    `json:"product_id" db:"product_id" validate:"required"`
	Description string   `json:"description" db:"description" validate:"required,min=1,max=500"`
	Rating      int      `json:"rating" db:"rating" validate:"required"`
	Images      []string `json:"images" db:"images" validate:"omitempty"`
}
