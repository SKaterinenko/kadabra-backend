package users_http

import (
	"context"
	user_model "kadabra/internal/features/users/model"
	users_service "kadabra/internal/features/users/service"
)

type UsersService interface {
	Register(ctx context.Context, req *users_service.CreateUserRequest) (*users_service.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*users_service.AuthResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*users_service.AuthResponse, error)
	GetByID(ctx context.Context, id int64) (*user_model.User, error)
}
