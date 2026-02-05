package users_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	user_model "kadabra/internal/features/users/model"
	"time"

	"kadabra/internal/core/config"
	"kadabra/pkg/utils"
)

type Service struct {
	repo     UserRepository
	s3Client *config.S3Client
	config   *config.Config
}

func NewService(repo UserRepository, s3Client *config.S3Client, cfg *config.Config) *Service {
	return &Service{repo: repo, s3Client: s3Client, config: cfg}
}

func (s *Service) Register(ctx context.Context, req *CreateUserRequest) (*AuthResponse, error) {
	if req.Password != req.RepeatPassword {
		return nil, errors.New("the passwords don't match")
	}

	// Проверяем, существует ли пользователь
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Если номер телефона указан, проверяем на уникальность
	if req.PhoneNumber != "" {
		existingPhone, err := s.repo.GetByPhoneNumber(ctx, req.PhoneNumber)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to check phone: %w", err)
		}
		if existingPhone != nil {
			return nil, errors.New("user with this phone number already exists")
		}
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

	// Подготавливаем phone_number (если пустой -> NULL)
	var phoneNumber *string
	if req.PhoneNumber != "" {
		phoneNumber = &req.PhoneNumber
	}

	// Создаём пользователя
	user := &user_model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		BirthDate:    birthDate,
		PhoneNumber:  phoneNumber, // nil если не указан
		Gender:       req.Gender,
		PasswordHash: passwordHash,
	}

	newUser, err := s.repo.Register(ctx, user)
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

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	// Получаем пользователя
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
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

func (s *Service) RefreshTokens(ctx context.Context, accessToken string) (*AuthResponse, error) {
	// Валидируем refresh token
	claims, err := utils.ValidateToken(accessToken, s.config.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Получаем пользователя
	user, err := s.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Генерируем новую пару токенов
	newAccessToken, newRefreshToken, err := utils.GenerateTokenPair(
		user.ID,
		user.Email,
		s.config.JWTSecret,
		s.config.JWTAccessExpiration,
		s.config.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*user_model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
