package categories_postgres

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core/config"
	categories_model "kadabra/internal/features/categories/model"
	categories_service "kadabra/internal/features/categories/service"
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
	var category categories_model.Category
	err = tx.QueryRow(
		ctx,
		"INSERT INTO categories DEFAULT VALUES RETURNING id, created_at, updated_at",
	).Scan(&category.Id, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert category: %w", err)
	}

	// 2. Подготавливаем вставку переводов
	insertBuilder := config.Psql.Insert("category_translations").
		Columns("category_id", "language_code", "name", "slug").
		Suffix("RETURNING id, category_id, language_code, name, slug, created_at, updated_at")

	// Сохраняем порядок языков для последующего чтения
	langCodes := make([]string, 0, len(req.Translations))

	for langCode, trans := range req.Translations {
		generatedSlug := slug.Make(trans.Name)

		insertBuilder = insertBuilder.Values(
			category.Id,
			langCode,
			trans.Name,
			generatedSlug,
		)

		langCodes = append(langCodes, langCode)
	}

	// 3. Вставляем переводы и читаем результаты
	insertTransSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build translations insert: %w", err)
	}

	rows, err := tx.Query(ctx, insertTransSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("insert translations: %w", err)
	}
	defer rows.Close()

	// 4. Читаем возвращенные данные с помощью pgx
	translationsList, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[categories_model.CategoryTranslation])
	if err != nil {
		return nil, fmt.Errorf("collect translations: %w", err)
	}

	// 5. Преобразуем список в map по language_code
	translations := make(map[string]categories_model.CategoryTranslation)
	for _, trans := range translationsList {
		translations[trans.LanguageCode] = *trans
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	result := &categories_model.CategoryWithTranslations{
		Id:           category.Id,
		Translations: translations,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}

	return result, nil
}

func (c *Category) GetAll(ctx context.Context, language string) ([]*categories_model.CategoryResponse, error) {
	query, args, err := config.Psql.
		Select(
			"c.id",
			"ct.name",
			"ct.slug",
			"c.created_at",
			"c.updated_at",
		).
		From("categories c").
		InnerJoin("category_translations ct ON c.id = ct.category_id").
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

	categories, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[categories_model.CategoryResponse])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	return categories, nil
}

func (c *Category) GetById(ctx context.Context, id int, language string) (*categories_model.CategoryResponse, error) {
	query, args, err := config.Psql.
		Select(
			"c.id",
			"ct.name",
			"ct.slug",
			"c.created_at",
			"c.updated_at",
		).
		From("categories c").
		InnerJoin("category_translations ct ON c.id = ct.category_id").
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

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.CategoryResponse])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found for language %s", language)
		}
		return nil, fmt.Errorf("collect row: %w", err)
	}

	return category, nil
}

func (c *Category) GetBySlug(ctx context.Context, slug, language string) (*categories_model.CategoryResponse, error) {
	query, args, err := config.Psql.
		Select(
			"c.id",
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

	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[categories_model.CategoryResponse])
	if err != nil {
		if err == pgx.ErrNoRows {
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
