package model

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Miner struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Energy    int       `json:"energy"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMiner(name string, energy, age int) *Miner {
	id := uuid.New()

	return &Miner{
		ID:        id,
		Name:      name,
		Energy:    energy,
		Age:       age,
		CreatedAt: time.Now(),
	}
}

type MinerRepository interface {
	Create(ctx context.Context, miner *Miner) error
	GetAll(ctx context.Context) ([]Miner, error)
}
