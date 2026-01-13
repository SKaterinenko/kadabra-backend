package products_type_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	products_type_model "kadabra/internal/features/products_type/model"
	products_type_service "kadabra/internal/features/products_type/service"
	sub_categories_model "kadabra/internal/features/sub_categories/model"
)

type ProductsType struct {
	db *pgxpool.Pool
}

func NewProductsTypePostgres(db *pgxpool.Pool) products_type_service.ProductsTypeRepository {
	return &ProductsType{db: db}
}

func (c *ProductsType) Create(ctx context.Context, productsType *products_type_model.ProductsType) (*products_type_model.ProductsType, error) {
	query, args, err := config.Psql.
		Insert("products_type").
		Columns("sub_category_id", "name").
		Values(productsType.SubCategoryId, productsType.Name).
		Suffix("RETURNING id, name, sub_category_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_type_model.ProductsType])
	if err != nil {
		return nil, core.QueryError(err)
	}

	return result, nil
}

func (c *ProductsType) GetAll(ctx context.Context) ([]*products_type_model.ProductsType, error) {
	query, _, err := config.Psql.
		Select("id", "name", "sub_category_id", "created_at", "updated_at").
		From("products_type").
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

	productsTypes, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_type_model.ProductsType])
	if err != nil {
		return nil, core.ScanError(err)
	}

	return productsTypes, nil
}

func (c *ProductsType) GetProductsTypeByCategorySlug(ctx context.Context, id string) ([]*sub_categories_model.SubCategoryWithProductsType, error) {
	query, args, err := config.Psql.
		Select(
			"sc.id",
			"sc.name",
			"sc.category_id",
			"sc.created_at",
			"sc.updated_at",
			"pt.id",
			"pt.name",
			"pt.sub_category_id",
			"pt.created_at",
			"pt.updated_at").
		From("sub_categories sc").
		Join("products_type pt ON sc.id = pt.sub_category_id").
		Where(sq.Eq{"sc.category_id": id}).
		Limit(30).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	mapSc := map[int]*sub_categories_model.SubCategoryWithProductsType{}

	for rows.Next() {
		var scWithProductType sub_categories_model.SubCategoryWithProductsType
		var productsType products_type_model.ProductsType

		if err := rows.Scan(
			&scWithProductType.Id,
			&scWithProductType.Name,
			&scWithProductType.CategoryId,
			&scWithProductType.CreatedAt,
			&scWithProductType.UpdatedAt,
			&productsType.Id,
			&productsType.Name,
			&productsType.SubCategoryId,
			&productsType.CreatedAt,
			&productsType.UpdatedAt,
		); err != nil {
			return nil, core.ScanError(err)
		}

		val, ok := mapSc[scWithProductType.Id]
		if !ok {
			scWithProductType.ProductsType = append(scWithProductType.ProductsType, &productsType)
			mapSc[scWithProductType.Id] = &scWithProductType
		} else {
			val.ProductsType = append(val.ProductsType, &productsType)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, core.RowsError(err)
	}

	subCategoryWithProductsTypes := make([]*sub_categories_model.SubCategoryWithProductsType, 0, len(mapSc))
	for _, val := range mapSc {
		subCategoryWithProductsTypes = append(subCategoryWithProductsTypes, val)
	}

	return subCategoryWithProductsTypes, nil
}

func (c *ProductsType) GetById(ctx context.Context, id int) (*products_type_model.ProductsType, error) {
	query, args, err := config.Psql.
		Select("id", "name", "sub_category_id", "created_at", "updated_at").
		From("products_type").
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

	productsType, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_type_model.ProductsType])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return productsType, nil
}

func (c *ProductsType) Delete(ctx context.Context, id int) error {
	query, args, err := config.Psql.Delete("products_type").Where(sq.Eq{"id": id}).ToSql()
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

func (c *ProductsType) Patch(ctx context.Context, id int, update *products_type_model.ProductsTypePatch) (*products_type_model.ProductsType, error) {
	q := config.Psql.
		Update("products_type").
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
		Suffix("RETURNING id, name, sub_category_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	productsType, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_type_model.ProductsType])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return productsType, nil
}
