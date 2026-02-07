package reviews_model

import "time"

type Review struct {
	ID          int64     `json:"id" db:"id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	ProductId   int64     `json:"product_id" db:"product_id"`
	Description string    `json:"description" db:"description"`
	Rating      int       `json:"rating" db:"rating"`
	Images      []string  `json:"images" db:"images"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ReviewWithUser struct {
	Review
	User ReviewUser `json:"user"`
}

type ReviewUser struct {
	ID        int64   `json:"id" db:"id"`
	FirstName string  `json:"first_name" db:"first_name"`
	LastName  string  `json:"last_name" db:"last_name"`
	Avatar    *string `json:"avatar,omitempty" db:"-"`
}

type Rating struct {
	TotalCount int `json:"total_count" db:"total_count"`
	Rating5    int `json:"rating_5"`
	Rating4    int `json:"rating_4"`
	Rating3    int `json:"rating_3"`
	Rating2    int `json:"rating_2"`
	Rating1    int `json:"rating_1"`
}

type ResReviews struct {
	Reviews []*ReviewWithUser `json:"reviews"`
	Ratings Rating            `json:"ratings"`
}
