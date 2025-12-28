package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"kadabra/internal/model"
	"kadabra/internal/service/manufacturerService"
)

type Manufacturer struct {
	db *sql.DB
}

func NewManufacturerPostgres(db *sql.DB) manufacturerService.ManufacturerRepository {
	return &Manufacturer{db: db}
}

func (m *Manufacturer) Create(ctx context.Context, manufacturer *model.Manufacturer) error {
	query, args, err := psql.
		Insert("manufacturers").
		Columns("id", "name").
		Values(manufacturer.Id, manufacturer.Name).
		Suffix("RETURNING created_at, updated_at").
		ToSql()

	if err != nil {
		return buildSQLError(err)
	}

	err = m.db.QueryRowContext(ctx, query, args...).Scan(&manufacturer.CreatedAt, &manufacturer.UpdatedAt)
	if err != nil {
		return queryError(err)
	}
	return nil
}

func (m *Manufacturer) GetAll(ctx context.Context) ([]*model.Manufacturer, error) {
	categories := []*model.Manufacturer{}
	query, _, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("manufacturers").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, rowsError(err)
	}

	for rows.Next() {
		var manufacturer model.Manufacturer
		if err := rows.Scan(&manufacturer.Id, &manufacturer.Name, &manufacturer.CreatedAt, &manufacturer.UpdatedAt); err != nil {
			return nil, scanError(err)
		}
		categories = append(categories, &manufacturer)
	}

	return categories, nil
}

func (m *Manufacturer) GetById(ctx context.Context, id uuid.UUID) (*model.Manufacturer, error) {
	var manufacturer model.Manufacturer
	query, args, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("manufacturers").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&manufacturer.Id, &manufacturer.Name, &manufacturer.CreatedAt, &manufacturer.UpdatedAt)
	if err != nil {
		return nil, NoGetRecord
	}
	return &manufacturer, nil
}

func (m *Manufacturer) Delete(ctx context.Context, id uuid.UUID) error {
	query, _, err := psql.Delete("manufacturers").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return buildSQLError(err)
	}
	result, err := m.db.ExecContext(ctx, query, id)
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

func (m *Manufacturer) Patch(ctx context.Context, id uuid.UUID, update *model.ManufacturerPatch) (*model.Manufacturer, error) {
	q := psql.
		Update("manufacturers").
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
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	var manufacturer model.Manufacturer
	err = m.db.QueryRowContext(ctx, query, args...).Scan(
		&manufacturer.Id,
		&manufacturer.Name,
		&manufacturer.CreatedAt,
		&manufacturer.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NoGetRecord
	}
	if err != nil {
		return nil, queryError(err)
	}

	return &manufacturer, nil
}
