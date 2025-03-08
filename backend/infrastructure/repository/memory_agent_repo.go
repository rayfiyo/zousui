package repository

import (
	"context"
	"errors"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type MemoryAgentRepo struct {
	Agents []*entity.Agent
}

func NewMemoryAgentRepo() *MemoryAgentRepo {
	return &MemoryAgentRepo{
		Agents: make([]*entity.Agent, 0),
	}
}

// ID に基づいてエージェントを返す
func (m *MemoryAgentRepo) GetByID(
	ctx context.Context,
	id string,
) (*entity.Agent, error) {
	logger := zap.L()
	logger.Debug("GetByID called", zap.String("agentID", id))
	for _, a := range m.Agents {
		if a.ID == id {
			return a, nil
		}
	}

	logger.Warn("Agent not found", zap.String("agentID", id))
	return nil, errors.New("agent not found")
}

// エージェントを保存する（既存なら更新、新規なら追加）
func (m *MemoryAgentRepo) Save(
	ctx context.Context, agent *entity.Agent,
) error {
	zap.L().Debug("Saving agent", zap.String("agentID", agent.ID))
	m.Agents = append(m.Agents, agent)
	zap.L().Info("Agent saved", zap.String("agentID", agent.ID))
	return nil
}

// communityID に基づくエージェントを返す
func (m *MemoryAgentRepo) GetAgentsByCommunity(
	ctx context.Context,
	communityID string,
) ([]*entity.Agent, error) {
	logger := zap.L()
	logger.Debug("GetAgentsByCommunity called", zap.String("communityID", communityID))

	// シンプルにフィルタ
	var result []*entity.Agent
	for _, a := range m.Agents {
		if a.CommunityID == communityID {
			result = append(result, a)
		}
	}
	logger.Info("Agents retrieved",
		zap.String("communityID", communityID), zap.Int("count", len(result)))
	return result, nil
}

// すべてのエージェントを返す
func (m *MemoryAgentRepo) GetAll(
	ctx context.Context,
) ([]*entity.Agent, error) {
	zap.L().Debug("GetAll agents called")
	return m.Agents, nil
}

var _ repository.AgentRepository = (*MemoryAgentRepo)(nil)
