package categories_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	"kadabra/internal/features/categories/model"
)

type Category struct {
	db *pgxpool.Pool
}

func NewCategoryPostgres(db *pgxpool.Pool) *Category {
	return &Category{db: db}
}

func (c *Category) Create(ctx context.Context, category *categories_model.Category) (*categories_model.Category, error) {
	query, args, err := config.Psql.
		Insert("categories").
		Columns("id", "name", "slug").
		Values(category.Id, category.Name, category.Slug).
		Suffix("RETURNING id, name, slug, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.Category])
	if err != nil {
		return nil, core.QueryError(err)
	}

	return result, nil
}

func (c *Category) GetAll(ctx context.Context) ([]*categories_model.Category, error) {
	query, _, err := config.Psql.
		Select("id", "name", "slug", "created_at", "updated_at").
		From("categories").
		Limit(30).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	categories, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[categories_model.Category])
	if err != nil {
		return nil, core.ScanError(err)
	}

	return categories, nil
}

func (c *Category) GetById(ctx context.Context, id uuid.UUID) (*categories_model.Category, error) {
	query, args, err := config.Psql.
		Select("id", "name", "slug", "created_at", "updated_at").
		From("categories").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return category, nil
}

func (c *Category) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := config.Psql.Delete("categories").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return core.BuildSQLError(err)
	}

	cmdTag, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return core.QueryError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return core.NoRecordDelete
	}
	return nil
}

func (c *Category) Patch(ctx context.Context, id uuid.UUID, update *categories_model.CategoryPatch) (*categories_model.Category, error) {
	q := config.Psql.
		Update("categories").
		Where(sq.Eq{"id": id})

	hasUpdates := false

	if update.Name != nil {
		q = q.Set("name", *update.Name)
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, core.NoUpdate
	}

	query, args, err := q.
		Suffix("RETURNING id, name, slug, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return category, nil
}
