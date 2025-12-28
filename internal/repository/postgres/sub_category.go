package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"kadabra/internal/model"
	"kadabra/internal/service/subCategoryService"
)

type SubCategory struct {
	db *sql.DB
}

func NewSubCategoryPostgres(db *sql.DB) subCategoryService.SubCategoryRepository {
	return &SubCategory{db: db}
}

func (c *SubCategory) Create(ctx context.Context, subCategory *model.SubCategory) error {
	query, args, err := psql.
		Insert("sub_categories").
		Columns("id", "category_id", "name").
		Values(subCategory.Id, subCategory.CategoryId, subCategory.Name).
		Suffix("RETURNING category_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return buildSQLError(err)
	}

	err = c.db.QueryRowContext(ctx, query, args...).Scan(&subCategory.CategoryId, &subCategory.CreatedAt, &subCategory.UpdatedAt)
	if err != nil {
		return queryError(err)
	}
	return nil
}

func (c *SubCategory) GetAll(ctx context.Context) ([]*model.SubCategory, error) {
	subCategories := []*model.SubCategory{}
	query, _, err := psql.
		Select("id", "name", "category_id", "created_at", "updated_at").
		From("sub_categories").
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
		var subCategory model.SubCategory
		if err := rows.Scan(&subCategory.Id, &subCategory.Name, &subCategory.CategoryId, &subCategory.CreatedAt, &subCategory.UpdatedAt); err != nil {
			return nil, scanError(err)
		}
		subCategories = append(subCategories, &subCategory)
	}

	return subCategories, nil
}

func (c *SubCategory) GetById(ctx context.Context, id uuid.UUID) (*model.SubCategory, error) {
	var subCategory model.SubCategory
	query, args, err := psql.
		Select("id", "name", "category_id", "created_at", "updated_at").
		From("sub_categories").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}
	err = c.db.QueryRowContext(ctx, query, args...).Scan(&subCategory.Id, &subCategory.Name, &subCategory.CategoryId, &subCategory.CreatedAt, &subCategory.UpdatedAt)
	if err != nil {
		return nil, NoGetRecord
	}
	return &subCategory, nil
}

func (c *SubCategory) Delete(ctx context.Context, id uuid.UUID) error {
	query, _, err := psql.Delete("sub_categories").Where(sq.Eq{"id": id}).ToSql()
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

func (c *SubCategory) Patch(ctx context.Context, id uuid.UUID, update *model.SubCategoryPatch) (*model.SubCategory, error) {
	q := psql.
		Update("sub_categories").
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
		Suffix("RETURNING id, name, category_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	var subCategory model.SubCategory
	err = c.db.QueryRowContext(ctx, query, args...).Scan(
		&subCategory.Id,
		&subCategory.Name,
		&subCategory.CategoryId,
		&subCategory.CreatedAt,
		&subCategory.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NoGetRecord
	}
	if err != nil {
		return nil, queryError(err)
	}

	return &subCategory, nil
}
