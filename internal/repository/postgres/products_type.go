package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/model"
	"kadabra/internal/service/productsTypeService"
)

type ProductsType struct {
	db *pgxpool.Pool
}

func NewProductsTypePostgres(db *pgxpool.Pool) productsTypeService.ProductsTypeRepository {
	return &ProductsType{db: db}
}

func (c *ProductsType) Create(ctx context.Context, productsType *model.ProductsType) (*model.ProductsType, error) {
	query, args, err := psql.
		Insert("products_type").
		Columns("id", "sub_category_id", "name").
		Values(productsType.Id, productsType.SubCategoryId, productsType.Name).
		Suffix("RETURNING id, name, sub_category_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.ProductsType])
	if err != nil {
		return nil, queryError(err)
	}

	return result, nil
}

func (c *ProductsType) GetAll(ctx context.Context) ([]*model.ProductsType, error) {
	query, _, err := psql.
		Select("id", "name", "sub_category_id", "created_at", "updated_at").
		From("products_type").
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

	productsTypes, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.ProductsType])
	if err != nil {
		return nil, scanError(err)
	}

	return productsTypes, nil
}

func (c *ProductsType) GetProductsTypeByCategoryId(ctx context.Context, id uuid.UUID) ([]*model.SubCategoryWithProductsType, error) {
	query, args, err := psql.
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
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	mapSc := map[uuid.UUID]*model.SubCategoryWithProductsType{}

	for rows.Next() {
		var scWithProductType model.SubCategoryWithProductsType
		var productsType model.ProductsType

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
			return nil, scanError(err)
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
		return nil, rowsError(err)
	}

	subCategoryWithProductsTypes := make([]*model.SubCategoryWithProductsType, 0, len(mapSc))
	for _, val := range mapSc {
		subCategoryWithProductsTypes = append(subCategoryWithProductsTypes, val)
	}

	return subCategoryWithProductsTypes, nil
}

func (c *ProductsType) GetById(ctx context.Context, id uuid.UUID) (*model.ProductsType, error) {
	query, args, err := psql.
		Select("id", "name", "sub_category_id", "created_at", "updated_at").
		From("products_type").
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

	productsType, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.ProductsType])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return productsType, nil
}

func (c *ProductsType) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("products_type").Where(sq.Eq{"id": id}).ToSql()
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

func (c *ProductsType) Patch(ctx context.Context, id uuid.UUID, update *model.SubCategoryPatch) (*model.ProductsType, error) {
	q := psql.
		Update("products_type").
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
		Suffix("RETURNING id, name, sub_category_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	productsType, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.ProductsType])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return productsType, nil
}
