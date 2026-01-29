package service

import (
	"context"
	"errors"
	user_model "kadabra/internal/features/users/model"
	"time"

	"kadabra/internal/core/config"
	"kadabra/pkg/utils"
)

type UserService struct {
	repo   UserRepository
	config *config.Config
}

func NewUserService(repo UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repo:   repo,
		config: cfg,
	}
}

func (s *UserService) Register(ctx context.Context, req *CreateUserRequest) (*AuthResponse, error) {
	// Проверяем, существует ли пользователь
	existingUser, _ := s.repo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хешируем пароль
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Парсим дату рождения
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, errors.New("invalid birth date format, use YYYY-MM-DD")
	}

	// Создаём пользователя
	user := &user_model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		BirthDate:    birthDate,
		PhoneNumber:  req.PhoneNumber,
		Gender:       req.Gender,
		PasswordHash: passwordHash,
	}

	newUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Генерируем токены
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		newUser.ID,
		newUser.Email,
		s.config.JWTSecret,
		s.config.JWTAccessExpiration,
		s.config.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         newUser,
	}, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	// Получаем пользователя
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Генерируем токены
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		user.ID,
		user.Email,
		s.config.JWTSecret,
		s.config.JWTAccessExpiration,
		s.config.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
