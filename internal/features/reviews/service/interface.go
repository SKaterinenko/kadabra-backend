package reviews_service

import (
	"context"
	reviews_model "kadabra/internal/features/reviews/model"
)

type ReviewsRepository interface {
	Create(ctx context.Context, review *CreateInput) (*reviews_model.Review, error)
	GetAllById(ctx context.Context, id, limit, offset int) (*reviews_model.ResReviews, error)
}
