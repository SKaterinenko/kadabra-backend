package service

import (
	"context"
	user_model "kadabra/internal/features/users/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *user_model.User) (*user_model.User, error)
	GetByEmail(ctx context.Context, email string) (*user_model.User, error)
	GetByID(ctx context.Context, id int64) (*user_model.User, error)
}
