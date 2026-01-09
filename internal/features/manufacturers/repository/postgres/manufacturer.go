package manufacturers_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	"kadabra/internal/features/manufacturers/model"
	manufacturers_service "kadabra/internal/features/manufacturers/service"
)

type Manufacturer struct {
	db *pgxpool.Pool
}

func NewManufacturerPostgres(db *pgxpool.Pool) manufacturers_service.ManufacturerRepository {
	return &Manufacturer{db: db}
}

func (m *Manufacturer) Create(ctx context.Context, manufacturer *manufacturers_model.Manufacturer) (*manufacturers_model.Manufacturer, error) {
	query, args, err := config.Psql.
		Insert("manufacturers").
		Columns("id", "name").
		Values(manufacturer.Id, manufacturer.Name).
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[manufacturers_model.Manufacturer])
	if err != nil {
		return nil, core.QueryError(err)
	}

	return result, nil
}

func (m *Manufacturer) GetAll(ctx context.Context) ([]*manufacturers_model.Manufacturer, error) {
	query, _, err := config.Psql.
		Select("id", "name", "created_at", "updated_at").
		From("manufacturers").
		Limit(30).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	manufacturers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[manufacturers_model.Manufacturer])
	if err != nil {
		return nil, core.ScanError(err)
	}

	return manufacturers, nil
}

func (m *Manufacturer) GetById(ctx context.Context, id uuid.UUID) (*manufacturers_model.Manufacturer, error) {
	query, args, err := config.Psql.
		Select("id", "name", "created_at", "updated_at").
		From("manufacturers").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	manufacturer, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[manufacturers_model.Manufacturer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return manufacturer, nil
}

func (m *Manufacturer) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := config.Psql.Delete("manufacturers").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return core.BuildSQLError(err)
	}

	cmdTag, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return core.QueryError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return core.NoRecordDelete
	}
	return nil
}

func (m *Manufacturer) Patch(ctx context.Context, id uuid.UUID, update *manufacturers_model.ManufacturerPatch) (*manufacturers_model.Manufacturer, error) {
	q := config.Psql.
		Update("manufacturers").
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
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	manufacturer, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[manufacturers_model.Manufacturer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return manufacturer, nil
}
