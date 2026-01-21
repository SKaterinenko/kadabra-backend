package products_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"

	"kadabra/internal/core"
	"kadabra/internal/core/config"
	products_model "kadabra/internal/features/products/model"
	products_service "kadabra/internal/features/products/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	db *pgxpool.Pool
}

func NewProductPostgres(db *pgxpool.Pool) *Product {
	return &Product{db: db}
}

func (c *Product) Create(ctx context.Context, req *products_service.CreateInput) (*products_model.ProductWithTranslations, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, args, err := config.Psql.
		Insert("products").
		Columns("product_type_id", "manufacturer_id").
		Values(req.ProductTypeId, req.ManufacturerId).
		Suffix("RETURNING id, product_type_id, manufacturer_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	productWT, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.ProductWithoutTranslations])
	if err != nil {
		return nil, core.QueryError(err)
	}

	product := &products_model.ProductWithTranslations{
		Id:             productWT.Id,
		ProductTypeId:  productWT.ProductTypeId,
		ManufacturerId: productWT.ManufacturerId,
		Translations:   make([]*products_model.ProductTranslate, 0),
		CreatedAt:      productWT.CreatedAt,
		UpdatedAt:      productWT.UpdatedAt,
	}

	query, args, err = config.Psql.
		Select("name").
		From("manufacturers").
		Where(sq.Eq{"id": req.ManufacturerId}).
		ToSql()

	rows, err = tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brand string

	for rows.Next() {
		err := rows.Scan(&brand)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range req.Translations {
		slugName := slug.Make(v.Name)
		slugText := brand + "-" + slugName

		query, args, err := config.Psql.
			Insert("product_translations").
			Columns("product_id", "language_code", "name", "slug", "short_description", "description").
			Values(productWT.Id, v.LanguageCode, v.Name, slugText, v.ShortDescription, v.Description).
			Suffix("RETURNING id, product_id, language_code, name, slug, short_description, description, created_at, updated_at").
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
			pgx.RowToAddrOfStructByName[products_model.ProductTranslate],
		)
		rows.Close()

		if err != nil {
			return nil, core.QueryError(err)
		}

		product.Translations = append(product.Translations, translate)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (c *Product) GetAll(ctx context.Context, lang string) ([]*products_model.Product, error) {
	query, args, err := config.Psql.
		Select("p.id", "pt.name", "pt.slug", "p.product_type_id", "p.manufacturer_id", "pt.short_description", "pt.description", "p.created_at", "p.updated_at").
		From("products p").
		Join("product_translations pt on p.id = pt.product_id").
		Where(sq.Eq{"pt.language_code": lang}).
		Limit(30).
		OrderBy("pt.name ASC").
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

func (c *Product) GetById(ctx context.Context, id int, lang string) (*products_model.Product, error) {
	query, args, err := config.Psql.
		Select(
			"p.id",
			"pt.name",
			"pt.slug",
			"p.product_type_id",
			"p.manufacturer_id",
			"pt.short_description",
			"pt.description",
			"p.created_at",
			"p.updated_at",
		).
		From("products p").
		Join("product_translations pt on p.id = pt.product_id").
		Where(sq.Eq{"pt.language_code": lang, "p.id": id}).
		OrderBy("pt.name ASC").
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
			product_type_id, 
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

func (c *Product) GetByCategoryIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error) {
	query, args, err := config.Psql.
		Select(
			"p.id",
			"pt.name",
			"pt.slug",
			"p.product_type_id",
			"p.manufacturer_id",
			"pt.short_description",
			"pt.description",
			"p.created_at",
			"p.updated_at",
		).
		From("products p").
		Join("product_translations pt on p.id = pt.product_id").
		Join("products_type pty ON p.product_type_id = pty.id").
		Join("sub_categories sc ON pty.sub_category_id = sc.id").
		Where(sq.Eq{"sc.category_id": categoryIds, "pt.language_code": lang}).
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

func (c *Product) GetByProductsTypeIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error) {
	sql, args, err := config.Psql.Select(
		"p.id",
		"pt.name",
		"pt.slug",
		"p.product_type_id",
		"p.manufacturer_id",
		"pt.short_description",
		"pt.description",
		"p.created_at",
		"p.updated_at").
		From("products p").
		Join("product_translations pt on p.id = pt.product_id").
		Where(sq.Eq{"p.product_type_id": categoryIds, "pt.language_code": lang}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c *Product) GetByCategorySlug(ctx context.Context, lang, slug string) ([]*products_model.Product, error) {
	sql, args, err := config.Psql.
		Select(
			"p.id",
			"pt.name",
			"pt.slug",
			"p.product_type_id",
			"p.manufacturer_id",
			"pt.short_description",
			"pt.description",
			"p.created_at",
			"p.updated_at",
		).From("products p").
		Join("product_translations pt on pt.product_id = p.id").
		Where(sq.Eq{"pt.language_code": lang}).
		Join("products_type pty on p.product_type_id = pty.id").
		Join("sub_categories sc on pty.sub_category_id = sc.id").
		Join("category_translations ct on sc.category_id = ct.category_id").
		Where(sq.Eq{"ct.slug": slug}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}
	rows, err := c.db.Query(ctx, sql, args...)
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

func (c *Product) GetByManufacturerId(ctx context.Context, id int, lang string) ([]*products_model.Product, error) {
	sql, args, err := config.Psql.Select(
		"p.id",
		"pt.name",
		"pt.slug",
		"p.product_type_id",
		"p.manufacturer_id",
		"pt.short_description",
		"pt.description",
		"p.created_at",
		"p.updated_at",
	).From("products p").
		Join("product_translations pt on pt.product_id = p.id").
		Where(sq.Eq{"pt.language_code": lang}).
		Join("manufacturers m on m.id = p.manufacturer_id").
		Where(sq.Eq{"m.id": id}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()
	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, err
	}

	return products, nil
}
