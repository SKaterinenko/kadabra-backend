package sub_categories_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
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

func NewSubCategoryPostgres(db *pgxpool.Pool) *SubCategory {
	return &SubCategory{db: db}
}

func (c *SubCategory) Create(ctx context.Context, req *sub_categories_service.CreateInput) (*sub_categories_model.SubCategoryWithTranslations, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, args, err := config.Psql.
		Insert("sub_categories").
		Columns("category_id").
		Values(req.CategoryId).
		Suffix("RETURNING id, category_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	subCategoryWT, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[sub_categories_model.SubCategoryWithoutTranslations])
	if err != nil {
		return nil, core.QueryError(err)
	}

	subCategory := &sub_categories_model.SubCategoryWithTranslations{
		Id:           subCategoryWT.Id,
		CategoryId:   subCategoryWT.CategoryId,
		Translations: make([]*sub_categories_model.SubCategoryTranslate, 0),
		CreatedAt:    subCategoryWT.CreatedAt,
		UpdatedAt:    subCategoryWT.UpdatedAt,
	}

	for _, v := range req.Translations {

		query, args, err := config.Psql.
			Insert("sub_category_translations").
			Columns("sub_category_id", "language_code", "name").
			Values(subCategoryWT.Id, v.LanguageCode, v.Name).
			Suffix("RETURNING id, sub_category_id, language_code, name, created_at, updated_at").
			ToSql()
		if err != nil {
			return nil, core.BuildSQLError(err)
		}

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, core.QueryError(err)
		}

		translate, err := pgx.CollectOneRow(
			rows,
			pgx.RowToAddrOfStructByName[sub_categories_model.SubCategoryTranslate],
		)
		rows.Close()

		if err != nil {
			return nil, core.QueryError(err)
		}

		subCategory.Translations = append(subCategory.Translations, translate)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return subCategory, nil
}

func (c *SubCategory) GetAll(ctx context.Context, lang string) ([]*sub_categories_model.SubCategory, error) {
	query, args, err := config.Psql.
		Select("sc.id", "sct.name", "sc.category_id", "sc.created_at", "sc.updated_at").
		From("sub_categories sc").
		Join("sub_category_translations sct on sc.id = sct.sub_category_id").
		Where(sq.Eq{"sct.language_code": lang}).
		Limit(30).
		OrderBy("sct.name ASC").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
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

func (c *SubCategory) GetById(ctx context.Context, id int, lang string) (*sub_categories_model.SubCategory, error) {
	query, args, err := config.Psql.
		Select("sc.id", "sct.name", "sc.category_id", "sc.created_at", "sc.updated_at").
		From("sub_categories sc").
		Join("sub_category_translations sct on sc.id = sct.sub_category_id").
		Where(sq.Eq{"sct.language_code": lang, "sc.id": id}).
		OrderBy("sct.name ASC").
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

func (c *SubCategory) Delete(ctx context.Context, id int) error {
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

func (c *SubCategory) Patch(ctx context.Context, id int, update *sub_categories_model.SubCategoryPatch) (*sub_categories_model.SubCategoryWithTranslations, error) {
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

	subCategory, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[sub_categories_model.SubCategoryWithTranslations])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return subCategory, nil
}
