package service

import (
	"context"
	"kadabra/internal/model"
)

type MinerService struct {
	repo model.MinerRepository
}

func NewMinerService(r model.MinerRepository) *MinerService {
	return &MinerService{repo: r}
}

func (m *MinerService) CreateMiner(ctx context.Context, miner *CreateMinerInput) (*CreateMinerOutput, error) {
	newMiner := model.NewMiner(miner.Name, miner.Energy, miner.Age)
	out := &CreateMinerOutput{
		ID:        newMiner.ID,
		Name:      newMiner.Name,
		Energy:    newMiner.Energy,
		Age:       newMiner.Age,
		CreatedAt: newMiner.CreatedAt,
	}
	if err := m.repo.Create(ctx, newMiner); err != nil {
		return nil, err
	}

	return out, nil
}

func (m *MinerService) GetAll(ctx context.Context) ([]model.Miner, error) {
	miners, err := m.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return miners, nil
}

// Множественный фильтр по query параметрам
// Просто добавить в конец
// if test != "" && test != v.Test {
//			continue
//		}

//func (c *MinerService) Filter(class string, work bool) []Miner {
//	result := []Miner{}
//	c.mtx.RLock()
//	defer c.mtx.RUnlock()
//
//	for _, v := range OurMiners {
//		if class != "" && class != v.Class {
//			continue
//		}
//		if work != nil && work != v.Work {
//			continue
//		}
//
//		result = append(result, v)
//	}
//	return result
//}
