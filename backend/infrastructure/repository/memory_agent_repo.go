package repository

import (
	"context"
	"errors"

	"github.com/rayfiyo/zousui/backend/domain"
	"github.com/rayfiyo/zousui/backend/usecase"
)

type MemoryAgentRepo struct {
	Agents []*domain.Agent
}

func NewMemoryAgentRepo() *MemoryAgentRepo {
	return &MemoryAgentRepo{
		Agents: make([]*domain.Agent, 0),
	}
}

// GetByID: ID に基づいてエージェントを返す
func (m *MemoryAgentRepo) GetByID(ctx context.Context, id string) (*domain.Agent, error) {
	for _, a := range m.Agents {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("agent not found")
}

// Save: エージェントを保存する（既存なら更新、新規なら追加）
func (m *MemoryAgentRepo) Save(ctx context.Context, agent *domain.Agent) error {
	m.Agents = append(m.Agents, agent)
	return nil
}

// GetAgentsByCommunity: communityID に基づくエージェントを返す
func (m *MemoryAgentRepo) GetAgentsByCommunity(ctx context.Context, communityID string) ([]*domain.Agent, error) {
	// シンプルにフィルタ
	var result []*domain.Agent
	for _, a := range m.Agents {
		if a.CommunityID == communityID {
			result = append(result, a)
		}
	}
	return result, nil
}

// GetAll: すべてのエージェントを返す
func (m *MemoryAgentRepo) GetAll(ctx context.Context) ([]*domain.Agent, error) {
	return m.Agents, nil
}

var _ usecase.AgentRepository = (*MemoryAgentRepo)(nil)
