package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"kadabra/internal/model"
	"kadabra/internal/service/productsTypeService"
)

type ProductsType struct {
	db *sql.DB
}

func NewProductsTypePostgres(db *sql.DB) productsTypeService.ProductsTypeRepository {
	return &ProductsType{db: db}
}

func (c *ProductsType) Create(ctx context.Context, productsType *model.ProductsType) error {
	query, args, err := psql.
		Insert("products_type").
		Columns("id", "sub_category_id", "name").
		Values(productsType.Id, productsType.SubCategoryId, productsType.Name).
		Suffix("RETURNING sub_category_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return buildSQLError(err)
	}

	err = c.db.QueryRowContext(ctx, query, args...).Scan(&productsType.SubCategoryId, &productsType.CreatedAt, &productsType.UpdatedAt)
	if err != nil {
		return queryError(err)
	}
	return nil
}

func (c *ProductsType) GetAll(ctx context.Context) ([]*model.ProductsType, error) {
	productsTypes := []*model.ProductsType{}
	query, _, err := psql.
		Select("id", "name", "sub_category_id", "created_at", "updated_at").
		From("products_type").
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
		var productsType model.ProductsType
		if err := rows.Scan(&productsType.Id, &productsType.Name, &productsType.SubCategoryId, &productsType.CreatedAt, &productsType.UpdatedAt); err != nil {
			return nil, scanError(err)
		}
		productsTypes = append(productsTypes, &productsType)
	}

	return productsTypes, nil
}

func (c *ProductsType) GetProductsTypeByCategoryId(ctx context.Context, id uuid.UUID) ([]*model.SubCategoryWithProductsType, error) {
	subCategoryWithProductsTypes := []*model.SubCategoryWithProductsType{}
	query, args, err := psql.
		Select("sc.id", "sc.name", "sc.category_id", "sc.created_at", "sc.updated_at", "pt.id", "pt.name", "pt.sub_category_id", "pt.created_at", "pt.updated_at").
		From("sub_categories sc").
		Join("products_type pt ON sc.id = pt.sub_category_id").
		Where(sq.Eq{"sc.category_id": id}).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, rowsError(err)
	}

	mapSc := map[uuid.UUID]*model.SubCategoryWithProductsType{}

	for rows.Next() {
		var scWithProductType model.SubCategoryWithProductsType
		var productsType model.ProductsType
		if err := rows.Scan(&scWithProductType.Id, &scWithProductType.Name, &scWithProductType.CategoryId, &scWithProductType.CreatedAt, &scWithProductType.UpdatedAt, &productsType.Id, &productsType.Name, &productsType.SubCategoryId, &productsType.CreatedAt, &productsType.UpdatedAt); err != nil {
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

	for _, val := range mapSc {
		subCategoryWithProductsTypes = append(subCategoryWithProductsTypes, val)
	}

	return subCategoryWithProductsTypes, nil
}

// select sc.name, pt.name  from sub_categories as sc join products_type as pt on sc.id = pt.sub_category_id where category_id = '057d7b5e-7821-4b2a-92e9-b1589c3228d5'
func (c *ProductsType) GetById(ctx context.Context, id uuid.UUID) (*model.ProductsType, error) {
	var productsType model.ProductsType
	query, args, err := psql.
		Select("id", "name", "sub_category_id", "created_at", "updated_at").
		From("products_type").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}
	err = c.db.QueryRowContext(ctx, query, args...).Scan(&productsType.Id, &productsType.Name, &productsType.SubCategoryId, &productsType.CreatedAt, &productsType.UpdatedAt)
	if err != nil {
		return nil, NoGetRecord
	}
	return &productsType, nil
}

func (c *ProductsType) Delete(ctx context.Context, id uuid.UUID) error {
	query, _, err := psql.Delete("products_type").Where(sq.Eq{"id": id}).ToSql()
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

	var productsType model.ProductsType
	err = c.db.QueryRowContext(ctx, query, args...).Scan(
		&productsType.Id,
		&productsType.Name,
		&productsType.SubCategoryId,
		&productsType.CreatedAt,
		&productsType.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NoGetRecord
	}
	if err != nil {
		return nil, queryError(err)
	}

	return &productsType, nil
}
