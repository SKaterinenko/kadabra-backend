package manufacturers_postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"

	"kadabra/internal/core"
	"kadabra/internal/core/config"
	manufacturers_model "kadabra/internal/features/manufacturers/model"
	manufacturers_service "kadabra/internal/features/manufacturers/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manufacturer struct {
	db *pgxpool.Pool
}

func NewManufacturerPostgres(db *pgxpool.Pool) manufacturers_service.ManufacturerRepository {
	return &Manufacturer{db: db}
}

func (m *Manufacturer) Create(ctx context.Context, req *manufacturers_service.CreateInput) (*manufacturers_model.ManufacturerWithTranslations, error) {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, args, err := config.Psql.
		Insert("manufacturers").
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	manufacturerWT, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[manufacturers_model.ManufacturerWithoutTranslations])
	if err != nil {
		return nil, core.QueryError(err)
	}

	manufacturer := &manufacturers_model.ManufacturerWithTranslations{
		Id:           manufacturerWT.Id,
		Translations: make([]*manufacturers_model.ManufacturerTranslate, 0),
		CreatedAt:    manufacturerWT.CreatedAt,
		UpdatedAt:    manufacturerWT.UpdatedAt,
	}

	for _, v := range req.Translations {
		slugText := slug.Make(v.Name)

		query, args, err := config.Psql.
			Insert("manufacturer_translations").
			Columns("manufacturer_id", "language_code", "name", "slug", "description").
			Values(manufacturerWT.Id, v.LanguageCode, v.Name, slugText, v.Description).
			Suffix("RETURNING id, manufacturer_id, language_code, name, slug, description, created_at, updated_at").
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
			pgx.RowToAddrOfStructByName[manufacturers_model.ManufacturerTranslate],
		)
		rows.Close()

		if err != nil {
			return nil, core.QueryError(err)
		}

		manufacturer.Translations = append(manufacturer.Translations, translate)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return manufacturer, nil
}

func (m *Manufacturer) GetAll(ctx context.Context, lang string) ([]*manufacturers_model.Manufacturer, error) {
	query, args, err := config.Psql.
		Select("m.id", "mt.name", "mt.slug", "mt.description", "m.created_at", "m.updated_at").
		From("manufacturers m").
		Join("manufacturer_translations mt on m.id = mt.manufacturer_id").
		Where(sq.Eq{"mt.language_code": lang}).
		Limit(30).
		OrderBy("mt.name ASC").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := m.db.Query(ctx, query, args...)
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

func (m *Manufacturer) GetById(ctx context.Context, id int, lang string) (*manufacturers_model.Manufacturer, error) {
	query, args, err := config.Psql.
		Select("m.id", "mt.name", "mt.slug", "mt.description", "m.created_at", "m.updated_at").
		From("manufacturers m").
		Join("manufacturer_translations mt on m.id = mt.manufacturer_id").
		Where(sq.Eq{"mt.language_code": lang, "m.id": id}).
		OrderBy("mt.name ASC").
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

func (m *Manufacturer) Delete(ctx context.Context, id int) error {
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

func (m *Manufacturer) Patch(ctx context.Context, id int, update *manufacturers_model.ManufacturerPatch) (*manufacturers_model.Manufacturer, error) {
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
		Suffix("RETURNING id, name, slug, created_at, updated_at").
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
