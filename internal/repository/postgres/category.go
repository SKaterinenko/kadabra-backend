package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"kadabra/internal/model"
	"kadabra/internal/service/categoryService"
)

type Category struct {
	db *sql.DB
}

func NewCategoryPostgres(db *sql.DB) categoryService.CategoryRepository {
	return &Category{db: db}
}

func (c *Category) Create(ctx context.Context, category *model.Category) error {
	query, args, err := psql.
		Insert("categories").
		Columns("id", "name").
		Values(category.Id, category.Name).
		Suffix("RETURNING created_at, updated_at").
		ToSql()

	if err != nil {
		return buildSQLError(err)
	}

	err = c.db.QueryRowContext(ctx, query, args...).Scan(&category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return queryError(err)
	}
	return nil
}

func (c *Category) GetAll(ctx context.Context) ([]*model.Category, error) {
	categories := []*model.Category{}
	query, _, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("categories").
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}
	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, rowsError(err)
	}

	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, scanError(err)
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (c *Category) GetById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	var category model.Category
	query, args, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("categories").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}
	err = c.db.QueryRowContext(ctx, query, args...).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, NoGetRecord
	}
	return &category, nil
}

func (c *Category) Delete(ctx context.Context, id uuid.UUID) error {
	query, _, err := psql.Delete("categories").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return buildSQLError(err)
	}
	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return queryError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return queryError(err)
	}
	if rows == 0 {
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

	var category model.Category
	err = c.db.QueryRowContext(ctx, query, args...).Scan(
		&category.Id,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NoGetRecord
	}
	if err != nil {
		return nil, queryError(err)
	}

	return &category, nil
}
