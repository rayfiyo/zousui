package repository

import (
	"context"

	"github.com/rayfiyo/zousui/backend/domain/entity"
)

type SimulationRepository interface {
	Save(ctx context.Context, result *entity.SimulationResult) error
	GetAll(ctx context.Context) ([]*entity.SimulationResult, error)
}
