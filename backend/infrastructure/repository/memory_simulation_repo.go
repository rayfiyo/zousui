package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rayfiyo/zousui/backend/domain/entity"
)

type MemorySimulationRepo struct {
	mu          sync.RWMutex
	simulations []*entity.SimulationResult
}

func NewMemorySimulationRepo() *MemorySimulationRepo {
	return &MemorySimulationRepo{
		simulations: make([]*entity.SimulationResult, 0),
	}
}

func (m *MemorySimulationRepo) Save(ctx context.Context, result *entity.SimulationResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// ID 自動生成（ここでは UUID を利用）
	result.ID = uuid.New().String()
	result.CreatedAt = time.Now()
	m.simulations = append(m.simulations, result)
	return nil
}

func (m *MemorySimulationRepo) GetAll(ctx context.Context) ([]*entity.SimulationResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.simulations, nil
}
