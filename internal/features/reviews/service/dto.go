package reviews_service

type CreateInput struct {
	UserId      int64    `json:"user_id" db:"user_id"`
	ProductId   int64    `json:"product_id" db:"product_id"`
	Description string   `json:"description" db:"description"`
	Rating      int      `json:"rating" db:"rating"`
	Images      []string `json:"images" db:"images"`
}
