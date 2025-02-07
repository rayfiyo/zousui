package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/rayfiyo/zousui/backend/domain"
	"github.com/rayfiyo/zousui/backend/usecase"
)

type MemoryCommunityRepo struct {
	mu          sync.RWMutex
	communities map[string]*domain.Community
}

func NewMemoryCommunityRepo() *MemoryCommunityRepo {
	return &MemoryCommunityRepo{
		communities: make(map[string]*domain.Community),
	}
}

func (m *MemoryCommunityRepo) GetByID(ctx context.Context, id string) (*domain.Community, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.communities[id]
	if !ok {
		return nil, errors.New("community not found")
	}
	return c, nil
}

func (m *MemoryCommunityRepo) Save(ctx context.Context, c *domain.Community) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.communities[c.ID] = c
	return nil
}

// 全コミュニティをリストとして取得
func (m *MemoryCommunityRepo) GetAll(ctx context.Context) ([]*domain.Community, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*domain.Community, 0, len(m.communities))
	for _, comm := range m.communities {
		result = append(result, comm)
	}
	return result, nil
}

// インタフェース実装をチェック
var _ usecase.CommunityRepository = (*MemoryCommunityRepo)(nil)
