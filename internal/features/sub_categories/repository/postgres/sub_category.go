package sub_categories_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	"kadabra/internal/features/sub_categories/model"
	"kadabra/internal/features/sub_categories/service"
)

type SubCategory struct {
	db *pgxpool.Pool
}

func NewSubCategoryPostgres(db *pgxpool.Pool) sub_categories_service.SubCategoryRepository {
	return &SubCategory{db: db}
}

func (c *SubCategory) Create(ctx context.Context, subCategory *sub_categories_model.SubCategory) (*sub_categories_model.SubCategory, error) {
	query, args, err := config.Psql.
		Insert("sub_categories").
		Columns("id", "category_id", "name").
		Values(subCategory.Id, subCategory.CategoryId, subCategory.Name).
		Suffix("RETURNING id, name, category_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[sub_categories_model.SubCategory])
	if err != nil {
		return nil, core.QueryError(err)
	}

	return result, nil
}

func (c *SubCategory) GetAll(ctx context.Context) ([]*sub_categories_model.SubCategory, error) {
	query, _, err := config.Psql.
		Select("id", "name", "category_id", "created_at", "updated_at").
		From("sub_categories").
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

	subCategories, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[sub_categories_model.SubCategory])
	if err != nil {
		return nil, core.ScanError(err)
	}

	return subCategories, nil
}

func (c *SubCategory) GetById(ctx context.Context, id uuid.UUID) (*sub_categories_model.SubCategory, error) {
	query, args, err := config.Psql.
		Select("id", "name", "category_id", "created_at", "updated_at").
		From("sub_categories").
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

	subCategory, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[sub_categories_model.SubCategory])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return subCategory, nil
}

func (c *SubCategory) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := config.Psql.Delete("sub_categories").Where(sq.Eq{"id": id}).ToSql()
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

func (c *SubCategory) Patch(ctx context.Context, id uuid.UUID, update *sub_categories_model.SubCategoryPatch) (*sub_categories_model.SubCategory, error) {
	q := config.Psql.
		Update("sub_categories").
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
		Suffix("RETURNING id, name, category_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	subCategory, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[sub_categories_model.SubCategory])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return subCategory, nil
}
