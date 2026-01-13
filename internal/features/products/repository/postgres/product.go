package products_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	"kadabra/internal/features/products/model"
)

type Product struct {
	db *pgxpool.Pool
}

func NewProductPostgres(db *pgxpool.Pool) *Product {
	return &Product{db: db}
}

func (c *Product) Create(ctx context.Context, product *products_model.Product) (*products_model.Product, error) {
	query, args, err := config.Psql.
		Insert("products").
		Columns("name", "products_type_id", "manufacturer_id", "short_description", "description").
		Values(
			product.Name,
			product.ProductsTypeId,
			product.ManufacturerId,
			product.ShortDescription,
			product.Description).
		Suffix("RETURNING id, name, products_type_id, manufacturer_id, short_description, description, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, core.QueryError(err)
	}

	return result, nil
}

func (c *Product) GetAll(ctx context.Context) ([]*products_model.Product, error) {
	query, _, err := config.Psql.
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
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, core.ScanError(err)
	}

	return products, nil
}

func (c *Product) GetById(ctx context.Context, id int) (*products_model.Product, error) {
	query, args, err := config.Psql.
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
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return product, nil
}

func (c *Product) Delete(ctx context.Context, id int) error {
	query, args, err := config.Psql.Delete("products").Where(sq.Eq{"id": id}).ToSql()
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

func (c *Product) Patch(ctx context.Context, id int, update *products_model.ProductPatch) (*products_model.Product, error) {
	q := config.Psql.
		Update("products").
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
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return product, nil
}

func (c *Product) GetByCategoryIds(ctx context.Context, categoryIds []int) ([]*products_model.Product, error) {
	if len(categoryIds) == 0 {
		return []*products_model.Product{}, nil
	}

	query, args, err := config.Psql.
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
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, core.ScanError(err)
	}

	return products, nil
}
