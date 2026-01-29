package service

import user_model "kadabra/internal/features/users/model"

type AuthResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	User         *user_model.User `json:"user"`
}

type CreateUserRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	BirthDate   string `json:"birth_date"` // "2000-01-15"
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	Password    string `json:"password"`
}
