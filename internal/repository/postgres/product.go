package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"kadabra/internal/model"
	"kadabra/internal/service/productService"
)

type Product struct {
	db *sql.DB
}

func NewProductPostgres(db *sql.DB) productService.ProductRepository {
	return &Product{db: db}
}

func (c *Product) Create(ctx context.Context, product *model.Product) error {
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
		Suffix("RETURNING created_at, updated_at").
		ToSql()

	if err != nil {
		return buildSQLError(err)
	}

	err = c.db.QueryRowContext(ctx, query, args...).Scan(&product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return queryError(err)
	}
	return nil
}

func (c *Product) GetAll(ctx context.Context) ([]*model.Product, error) {
	products := []*model.Product{}
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
		var product model.Product
		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.ProductsTypeId,
			&product.ManufacturerId,
			&product.ShortDescription,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt); err != nil {
			return nil, scanError(err)
		}
		products = append(products, &product)
	}
	return products, nil
}

func (c *Product) GetById(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var product model.Product
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
	err = c.db.QueryRowContext(ctx, query, args...).Scan(
		&product.Id,
		&product.Name,
		&product.ProductsTypeId,
		&product.ManufacturerId,
		&product.ShortDescription,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, NoGetRecord
	}
	return &product, nil
}

func (c *Product) Delete(ctx context.Context, id uuid.UUID) error {
	query, _, err := psql.Delete("products").Where(sq.Eq{"id": id}).ToSql()
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
		Suffix("RETURNING id, name, products_type_id, manufacturer_id, short_description, description  created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	var product model.Product
	err = c.db.QueryRowContext(ctx, query, args...).Scan(
		&product.Id,
		&product.Name,
		&product.ProductsTypeId,
		&product.ManufacturerId,
		&product.ShortDescription,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NoGetRecord
	}
	if err != nil {
		return nil, queryError(err)
	}

	return &product, nil
}

func (c *Product) GetByCategoryIds(ctx context.Context, categoryIds []uuid.UUID) ([]*model.Product, error) {
	products := []*model.Product{}

	if len(categoryIds) == 0 {
		return products, nil
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
		OrderBy("created_at ASC").
		ToSql()

	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.ProductsTypeId,
			&product.ManufacturerId,
			&product.ShortDescription,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
