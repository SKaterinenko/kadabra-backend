package reviews_service

import (
	"context"
	"kadabra/internal/core/config"
	reviews_model "kadabra/internal/features/reviews/model"
)

type Service struct {
	repo     ReviewsRepository
	s3Client *config.S3Client
	config   *config.Config
}

func NewService(repo ReviewsRepository, s3Client *config.S3Client, cfg *config.Config) *Service {
	return &Service{repo: repo, s3Client: s3Client, config: cfg}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (*reviews_model.Review, error) {
	if input.Images == nil {
		input.Images = []string{}
	}
	review, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (s *Service) GetAllById(ctx context.Context, id, limit, offset int) (*reviews_model.ResReviews, error) {
	if limit == 0 {
		limit = 20
	}
	reviews, err := s.repo.GetAllById(ctx, id, limit, offset)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
