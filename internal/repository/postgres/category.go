package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/model"
	"kadabra/internal/service/categoryService"
)

type Category struct {
	db *pgxpool.Pool
}

func NewCategoryPostgres(db *pgxpool.Pool) categoryService.CategoryRepository {
	return &Category{db: db}
}

func (c *Category) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	query, args, err := psql.
		Insert("categories").
		Columns("id", "name").
		Values(category.Id, category.Name).
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Category])
	if err != nil {
		return nil, queryError(err)
	}

	return result, nil
}

func (c *Category) GetAll(ctx context.Context) ([]*model.Category, error) {
	query, _, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("categories").
		Limit(30).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	categories, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Category])
	if err != nil {
		return nil, scanError(err)
	}

	return categories, nil
}

func (c *Category) GetById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	query, args, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("categories").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return category, nil
}

func (c *Category) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("categories").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return buildSQLError(err)
	}

	cmdTag, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return queryError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return NoRecordDelete
	}
	return nil
}

func (c *Category) Patch(ctx context.Context, id uuid.UUID, update *model.CategoryPatch) (*model.Category, error) {
	q := psql.
		Update("categories").
		Where(sq.Eq{"id": id})

	hasUpdates := false

	if update.Name != nil {
		q = q.Set("name", *update.Name)
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, NoUpdate
	}

	query, args, err := q.
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return category, nil
}
