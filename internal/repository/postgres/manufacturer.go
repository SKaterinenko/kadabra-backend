package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/model"
	"kadabra/internal/service/manufacturerService"
)

type Manufacturer struct {
	db *pgxpool.Pool
}

func NewManufacturerPostgres(db *pgxpool.Pool) manufacturerService.ManufacturerRepository {
	return &Manufacturer{db: db}
}

func (m *Manufacturer) Create(ctx context.Context, manufacturer *model.Manufacturer) (*model.Manufacturer, error) {
	query, args, err := psql.
		Insert("manufacturers").
		Columns("id", "name").
		Values(manufacturer.Id, manufacturer.Name).
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Manufacturer])
	if err != nil {
		return nil, queryError(err)
	}

	return result, nil
}

func (m *Manufacturer) GetAll(ctx context.Context) ([]*model.Manufacturer, error) {
	query, _, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("manufacturers").
		Limit(30).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	manufacturers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Manufacturer])
	if err != nil {
		return nil, scanError(err)
	}

	return manufacturers, nil
}

func (m *Manufacturer) GetById(ctx context.Context, id uuid.UUID) (*model.Manufacturer, error) {
	query, args, err := psql.
		Select("id", "name", "created_at", "updated_at").
		From("manufacturers").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	manufacturer, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Manufacturer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return manufacturer, nil
}

func (m *Manufacturer) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("manufacturers").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return buildSQLError(err)
	}

	cmdTag, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return queryError(err)
	}

	if cmdTag.RowsAffected() == 0 {
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

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	manufacturer, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Manufacturer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoGetRecord
		}
		return nil, queryError(err)
	}

	return manufacturer, nil
}
