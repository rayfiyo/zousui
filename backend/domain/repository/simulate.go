package repository

import (
	"context"

	"github.com/rayfiyo/zousui/backend/domain/entity"
)

// AgentRepository: エージェントに関するリポジトリインタフェース(読み書き)
type AgentRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Agent, error)
	Save(ctx context.Context, agent *entity.Agent) error
	GetAll(ctx context.Context) ([]*entity.Agent, error)
	GetAgentsByCommunity(ctx context.Context, communityID string) ([]*entity.Agent, error)
}

// CommunityRepository: コミュニティに関するリポジトリインタフェース(読み書き)
type CommunityRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Community, error)
	Save(ctx context.Context, community *entity.Community) error
	GetAll(ctx context.Context) ([]*entity.Community, error)
	Delete(ctx context.Context, id string) error
}

// LLMGateway: LLMに問い合わせるためのインタフェース
type LLMGateway interface {
	GenerateCultureUpdate(ctx context.Context, prompt string) (string, error)
}
