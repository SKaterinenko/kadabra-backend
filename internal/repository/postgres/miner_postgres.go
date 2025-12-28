package repository

import (
	"context"
	"database/sql"
	"fmt"
	"kadabra/internal/model"
)

type MinerPostgres struct {
	db *sql.DB
}

func NewMinerPostgres(db *sql.DB) model.MinerRepository {
	return &MinerPostgres{db: db}
}

//db.ExecContext - это insert/update/delete без возрата данных
//db.QueryRowContext select/insert/update/delete с вовзратом 1 строки если есть RETURNING
//db.QueryContext select возвращает много строк. редкий случай для insert

//func (r *MinerPostgres) Create1(ctx context.Context, m *model.Miner) (int64, error) {
//	var id int64
//	err := r.db.QueryRowContext(ctx,
//		`INSERT INTO miners(name, energy) VALUES($1, $2) RETURNING id`,
//		m.Name, m.Energy).Scan(&id)
//	return id, err
//}

func (r *MinerPostgres) Create(ctx context.Context, miner *model.Miner) error {
	query := `INSERT INTO miners (id, name, energy, age) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, miner.ID, miner.Name, miner.Energy, miner.Age)
	if err != nil {
		return fmt.Errorf("failed to create miner: %w", err)
	}
	return nil
}

func (r *MinerPostgres) GetAll(ctx context.Context) ([]model.Miner, error) {
	miners := []model.Miner{}
	query := `SELECT id, name, energy, age, created_at FROM miners`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get miners: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var miner model.Miner
		if err := rows.Scan(&miner.ID, &miner.Name, &miner.Energy, &miner.Age, &miner.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan miner: %w", err)
		}
		miners = append(miners, miner)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return miners, nil
}
