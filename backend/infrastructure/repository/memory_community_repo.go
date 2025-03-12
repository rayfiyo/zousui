package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type MemoryCommunityRepo struct {
	mu          sync.RWMutex
	communities map[string]*entity.Community
}

func NewMemoryCommunityRepo() *MemoryCommunityRepo {
	zap.L().Debug("Initializing MemoryCommunityRepo")
	return &MemoryCommunityRepo{
		communities: make(map[string]*entity.Community),
	}
}

// IDでコミュニティを取得
func (m *MemoryCommunityRepo) GetByID(
	ctx context.Context,
	id string,
) (*entity.Community, error) {
	logger := zap.L()
	logger.Debug("GetByID called", zap.String("communityID", id))
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.communities[id]
	if !ok {
		logger.Warn("Community not found", zap.String("communityID", id))
		return nil, errors.New("community not found")
	}
	logger.Debug("Community found", zap.String("communityID", id))
	return c, nil
}

// コミュニティを保存
func (m *MemoryCommunityRepo) Save(
	ctx context.Context,
	c *entity.Community,
) error {
	logger := zap.L()
	logger.Debug("Saving community", zap.String("communityID", c.ID))
	m.mu.Lock()
	defer m.mu.Unlock()
	m.communities[c.ID] = c
	logger.Info("Community saved", zap.String("communityID", c.ID))
	return nil
}

// 全コミュニティをリストとして取得
func (m *MemoryCommunityRepo) GetAll(
	ctx context.Context,
) ([]*entity.Community, error) {
	logger := zap.L()
	logger.Debug("GetAll communities called")
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*entity.Community, 0, len(m.communities))
	for _, comm := range m.communities {
		result = append(result, comm)
	}
	logger.Info("Retrieved all communities", zap.Int("count", len(result)))
	return result, nil
}

// コミュニティを削除
func (m *MemoryCommunityRepo) Delete(
	ctx context.Context,
	id string,
) error {
	logger := zap.L()
	logger.Debug("Delete called", zap.String("communityID", id))
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.communities[id]; !ok {
		logger.Warn("Community to delete not found", zap.String("communityID", id))
		return errors.New("community not found")
	}
	delete(m.communities, id)
	logger.Info("Community deleted", zap.String("communityID", id))
	return nil
}

// インタフェース実装をチェック
var _ repository.CommunityRepository = (*MemoryCommunityRepo)(nil)
