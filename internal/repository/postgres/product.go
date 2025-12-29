package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/model"
	"kadabra/internal/service/productService"
)

type Product struct {
	db *pgxpool.Pool
}

func NewProductPostgres(db *pgxpool.Pool) productService.ProductRepository {
	return &Product{db: db}
}

func (c *Product) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	query, args, err := psql.
		Insert("products").
		Columns("id", "name", "products_type_id", "manufacturer_id", "short_description", "description").
		Values(
			product.Id,
			product.Name,
			product.ProductsTypeId,
			product.ManufacturerId,
			product.ShortDescription,
			product.Description).
		Suffix("RETURNING id, name, products_type_id, manufacturer_id, short_description, description, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Product])
	if err != nil {
		return nil, queryError(err)
	}

	return result, nil
}

func (c *Product) GetAll(ctx context.Context) ([]*model.Product, error) {
	query, _, err := psql.
		Select(
			"id",
			"name",
			"products_type_id",
			"manufacturer_id",
			"short_description",
			"description",
			"created_at",
			"updated_at").
		From("products").
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

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Product])
	if err != nil {
		return nil, scanError(err)
	}

	return products, nil
}

func (c *Product) GetById(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	query, args, err := psql.
		Select(
			"id",
			"name",
			"products_type_id",
			"manufacturer_id",
			"short_description",
			"description",
			"created_at",
			"updated_at",
		).
		From("products").
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

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Product])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return product, nil
}

func (c *Product) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("products").Where(sq.Eq{"id": id}).ToSql()
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

func (c *Product) Patch(ctx context.Context, id uuid.UUID, update *model.ProductPatch) (*model.Product, error) {
	q := psql.
		Update("products").
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
		Suffix(`
			RETURNING id,
			name, 
			products_type_id, 
			manufacturer_id, 
			short_description, 
			description, 
			created_at, 
			updated_at`).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Product])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return product, nil
}

func (c *Product) GetByCategoryIds(ctx context.Context, categoryIds []uuid.UUID) ([]*model.Product, error) {
	if len(categoryIds) == 0 {
		return []*model.Product{}, nil
	}

	query, args, err := psql.
		Select(
			"p.id",
			"p.name",
			"p.products_type_id",
			"p.manufacturer_id",
			"p.short_description",
			"p.description",
			"p.created_at",
			"p.updated_at",
		).
		From("products p").
		Join("products_type pt ON p.products_type_id = pt.id").
		Join("sub_categories sc ON pt.sub_category_id = sc.id").
		Where(sq.Eq{"sc.category_id": categoryIds}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()

	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Product])
	if err != nil {
		return nil, scanError(err)
	}

	return products, nil
}
