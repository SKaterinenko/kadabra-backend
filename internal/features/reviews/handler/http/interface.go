package reviews_http

import (
	"context"
	reviews_model "kadabra/internal/features/reviews/model"
	reviews_service "kadabra/internal/features/reviews/service"
)

type ReviewsService interface {
	Create(ctx context.Context, review *reviews_service.CreateInput) (*reviews_model.Review, error)
	GetAllById(ctx context.Context, id, limit, offset int) (*reviews_model.ResReviews, error)
}
