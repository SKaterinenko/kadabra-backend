package categories_postgres

import (
	"context"
	"errors"
	"fmt"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	categories_model "kadabra/internal/features/categories/model"
	categories_service "kadabra/internal/features/categories/service"

	sq "github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	db *pgxpool.Pool
}

func NewCategoryPostgres(db *pgxpool.Pool) *Category {
	return &Category{db: db}
}

func (c *Category) Create(ctx context.Context, req *categories_service.CreateInput) (*categories_model.CategoryWithTranslations, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Создаем категорию
	var categoryWT categories_model.Category
	err = tx.QueryRow(
		ctx,
		"INSERT INTO categories DEFAULT VALUES RETURNING id, created_at, updated_at",
	).Scan(&categoryWT.Id, &categoryWT.CreatedAt, &categoryWT.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert category: %w", err)
	}

	category := &categories_model.CategoryWithTranslations{
		Id:           categoryWT.Id,
		Translations: make([]*categories_model.CategoryTranslate, 0),
		CreatedAt:    categoryWT.CreatedAt,
		UpdatedAt:    categoryWT.UpdatedAt,
	}

	for _, v := range req.Translations {
		generatedSlug := slug.Make(v.Name)

		query, args, err := config.Psql.
			Insert("category_translations").
			Columns("category_id", "language_code", "name", "slug").
			Values(categoryWT.Id, v.LanguageCode, v.Name, generatedSlug).
			Suffix("RETURNING id, category_id, language_code, name, slug, created_at, updated_at").
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("build translation insert: %w", err)
		}

		rows, err := tx.Query(ctx, query, args...)

		if err != nil {
			return nil, fmt.Errorf("insert translation: %w", err)
		}

		translate, err := pgx.CollectOneRow(
			rows,
			pgx.RowToAddrOfStructByName[categories_model.CategoryTranslate],
		)
		rows.Close()

		if err != nil {
			return nil, fmt.Errorf("collect translation: %w", err)
		}

		category.Translations = append(category.Translations, translate)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return category, nil
}

func (c *Category) GetAll(ctx context.Context, language string) ([]*categories_model.Category, error) {
	query, args, err := config.Psql.
		Select(
			"c.id",
			"c.image",
			"ct.name",
			"ct.slug",
			"c.created_at",
			"c.updated_at",
		).
		From("categories c").
		Join("category_translations ct ON c.id = ct.category_id").
		Where(sq.Eq{"ct.language_code": language}).
		OrderBy("c.id ASC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	categories, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[categories_model.Category])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	return categories, nil
}

func (c *Category) GetById(ctx context.Context, id int, language string) (*categories_model.Category, error) {
	query, args, err := config.Psql.
		Select(
			"c.id",
			"c.image",
			"ct.name",
			"ct.slug",
			"c.created_at",
			"c.updated_at",
		).
		From("categories c").
		Join("category_translations ct ON c.id = ct.category_id").
		Where(sq.Eq{
			"c.id":             id,
			"ct.language_code": language,
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query category: %w", err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found for language %s", language)
		}
		return nil, fmt.Errorf("collect row: %w", err)
	}

	return category, nil
}

func (c *Category) GetBySlug(ctx context.Context, slug, language string) (*categories_model.Category, error) {
	query, args, err := config.Psql.
		Select(
			"c.id",
			"c.image",
			"ct.name",
			"ct.slug",
			"c.created_at",
			"c.updated_at",
		).
		From("categories c").
		InnerJoin("category_translations ct ON c.id = ct.category_id").
		Where(sq.Eq{
			"ct.slug":          slug,
			"ct.language_code": language,
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query category: %w", err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("collect row: %w", err)
	}

	return category, nil
}

func (c *Category) Delete(ctx context.Context, id int) error {
	sql, args, err := config.Psql.
		Delete("categories").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	commandTag, err := c.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("category with id %d not found", id)
	}

	return nil
}

func (c *Category) Patch(ctx context.Context, id int, update *categories_model.CategoryPatch) (*categories_model.CategoryWithoutTranslations, error) {
	q := config.Psql.
		Update("categories").
		Where(sq.Eq{"id": id})

	hasUpdates := false

	if update.Image != nil {
		q = q.Set("image", *update.Image)
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, core.NoUpdate
	}

	query, args, err := q.
		Suffix("RETURNING id, image, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.CategoryWithoutTranslations])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return category, nil
}
