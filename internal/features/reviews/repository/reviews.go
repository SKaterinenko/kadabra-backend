package reviews_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	reviews_model "kadabra/internal/features/reviews/model"
	reviews_service "kadabra/internal/features/reviews/service"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewReviewsPostgres(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, review *reviews_service.CreateReq) (*reviews_model.Review, error) {
	query, args, err := config.Psql.
		Insert("reviews").
		Columns(
			"product_id",
			"user_id",
			"rating",
			"description",
			"images",
		).
		Values(
			review.ProductId,
			review.UserId,
			review.Rating,
			review.Description,
			review.Images,
		).
		Suffix("RETURNING id, product_id, user_id, rating, description, images, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	newReview, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[reviews_model.Review])
	if err != nil {
		return nil, core.RowsError(err)
	}

	return newReview, nil
}

func (r *Repository) GetAllById(ctx context.Context, id, limit, offset int) (*reviews_model.ResReviews, error) {
	sql, args, err := config.Psql.
		Select(
			"r.id",
			"r.product_id",
			"r.user_id",
			"r.rating",
			"r.description",
			"r.images",
			"u.id",
			"u.first_name",
			"u.last_name",
			"u.avatar",
			"r.created_at",
			"r.updated_at").
		From("reviews r").
		Where(sq.Eq{"product_id": id}).
		Join("users u on r.user_id = u.id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()
	var reviews reviews_model.ResReviews
	for rows.Next() {
		var review reviews_model.ReviewWithUser
		err := rows.Scan(
			&review.Review.ID,
			&review.Review.ProductId,
			&review.Review.UserId,
			&review.Review.Rating,
			&review.Review.Description,
			&review.Review.Images,
			&review.User.ID,
			&review.User.FirstName,
			&review.User.LastName,
			&review.User.Avatar,
			&review.Review.CreatedAt,
			&review.Review.UpdatedAt,
		)
		if err != nil {
			return nil, core.ScanError(err)
		}
		reviews.Reviews = append(reviews.Reviews, &review)
	}
	for i := 1; i <= 5; i++ {
		sql, args, err = config.Psql.
			Select("COUNT(*)").
			From("reviews").
			Where(sq.Eq{
				"product_id": id,
				"rating":     i,
			}).
			ToSql()
		if err != nil {
			return nil, core.BuildSQLError(err)
		}
		var count int
		err = r.db.QueryRow(ctx, sql, args...).Scan(&count)
		if err != nil {
			return nil, core.QueryError(err)
		}

		switch i {
		case 1:
			reviews.Ratings.Rating1 = count
		case 2:
			reviews.Ratings.Rating2 = count
		case 3:
			reviews.Ratings.Rating3 = count
		case 4:
			reviews.Ratings.Rating4 = count
		case 5:
			reviews.Ratings.Rating5 = count
		}

		sql, args, err = config.Psql.
			Select("COUNT(*)").
			From("reviews").
			Where(sq.Eq{
				"product_id": id,
			}).
			ToSql()
		if err != nil {
			return nil, core.BuildSQLError(err)
		}

		err = r.db.QueryRow(ctx, sql, args...).Scan(&reviews.Ratings.TotalCount)
		if err != nil {
			return nil, core.QueryError(err)
		}

	}
	return &reviews, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	sql, args, err := config.Psql.Delete("reviews").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return core.BuildSQLError(err)
	}
	cmd, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return core.QueryError(err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("review not found")
	}
	return nil
}
