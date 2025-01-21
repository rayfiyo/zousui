package repository

import (
	"context"
	"errors"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
)

type MemoryAgentRepo struct {
	Agents []*entity.Agent
}

func NewMemoryAgentRepo() *MemoryAgentRepo {
	return &MemoryAgentRepo{
		Agents: make([]*entity.Agent, 0),
	}
}

// GetByID: ID に基づいてエージェントを返す
func (m *MemoryAgentRepo) GetByID(ctx context.Context, id string) (*entity.Agent, error) {
	for _, a := range m.Agents {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("agent not found")
}

// Save: エージェントを保存する（既存なら更新、新規なら追加）
func (m *MemoryAgentRepo) Save(ctx context.Context, agent *entity.Agent) error {
	m.Agents = append(m.Agents, agent)
	return nil
}

// GetAgentsByCommunity: communityID に基づくエージェントを返す
func (m *MemoryAgentRepo) GetAgentsByCommunity(ctx context.Context, communityID string) ([]*entity.Agent, error) {
	// シンプルにフィルタ
	var result []*entity.Agent
	for _, a := range m.Agents {
		if a.CommunityID == communityID {
			result = append(result, a)
		}
	}
	return result, nil
}

// GetAll: すべてのエージェントを返す
func (m *MemoryAgentRepo) GetAll(ctx context.Context) ([]*entity.Agent, error) {
	return m.Agents, nil
}

var _ repository.AgentRepository = (*MemoryAgentRepo)(nil)
