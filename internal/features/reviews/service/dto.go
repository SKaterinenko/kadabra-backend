package reviews_service

import "mime/multipart"

type CreateInput struct {
	UserId      int64                   `json:"user_id" db:"user_id"`
	ProductId   int                     `json:"product_id" db:"product_id"`
	Description string                  `json:"description" db:"description"`
	Rating      int                     `json:"rating" db:"rating"`
	Images      []*multipart.FileHeader `json:"images" db:"images"`
}

type CreateReq struct {
	UserId      int64    `json:"user_id" db:"user_id"`
	ProductId   int      `json:"product_id" db:"product_id"`
	Description string   `json:"description" db:"description"`
	Rating      int      `json:"rating" db:"rating"`
	Images      []string `json:"images" db:"images"`
}
