package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/model"
	"kadabra/internal/service/subCategoryService"
)

type SubCategory struct {
	db *pgxpool.Pool
}

func NewSubCategoryPostgres(db *pgxpool.Pool) subCategoryService.SubCategoryRepository {
	return &SubCategory{db: db}
}

func (c *SubCategory) Create(ctx context.Context, subCategory *model.SubCategory) (*model.SubCategory, error) {
	query, args, err := psql.
		Insert("sub_categories").
		Columns("id", "category_id", "name").
		Values(subCategory.Id, subCategory.CategoryId, subCategory.Name).
		Suffix("RETURNING id, name, category_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.SubCategory])
	if err != nil {
		return nil, queryError(err)
	}

	return result, nil
}

func (c *SubCategory) GetAll(ctx context.Context) ([]*model.SubCategory, error) {
	query, _, err := psql.
		Select("id", "name", "category_id", "created_at", "updated_at").
		From("sub_categories").
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

	subCategories, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.SubCategory])
	if err != nil {
		return nil, scanError(err)
	}

	return subCategories, nil
}

func (c *SubCategory) GetById(ctx context.Context, id uuid.UUID) (*model.SubCategory, error) {
	query, args, err := psql.
		Select("id", "name", "category_id", "created_at", "updated_at").
		From("sub_categories").
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

	subCategory, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.SubCategory])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return subCategory, nil
}

func (c *SubCategory) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("sub_categories").Where(sq.Eq{"id": id}).ToSql()
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

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	subCategory, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.SubCategory])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return subCategory, nil
}
