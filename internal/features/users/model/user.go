package users_model

import "time"

type User struct {
	ID           int64     `json:"id" db:"id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email" db:"email"`
	BirthDate    time.Time `json:"birth_date" db:"birth_date"`
	PhoneNumber  *string   `json:"phone_number,omitempty" db:"phone_number"`
	Gender       string    `json:"gender" db:"gender"`
	Avatar       *string   `json:"avatar,omitempty" db:"-"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
