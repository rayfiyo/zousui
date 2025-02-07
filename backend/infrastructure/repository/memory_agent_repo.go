package repository

import (
	"context"

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

var _ usecase.AgentRepository = (*MemoryAgentRepo)(nil)
