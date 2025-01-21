package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
)

type MemoryCommunityRepo struct {
	mu          sync.RWMutex
	communities map[string]*entity.Community
}

func NewMemoryCommunityRepo() *MemoryCommunityRepo {
	return &MemoryCommunityRepo{
		communities: make(map[string]*entity.Community),
	}
}

// GetByID: IDでコミュニティを取得
func (m *MemoryCommunityRepo) GetByID(ctx context.Context, id string) (*entity.Community, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.communities[id]
	if !ok {
		return nil, errors.New("community not found")
	}
	return c, nil
}

// Save: コミュニティを保存
func (m *MemoryCommunityRepo) Save(ctx context.Context, c *entity.Community) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.communities[c.ID] = c
	return nil
}

// Get All: 全コミュニティをリストとして取得
func (m *MemoryCommunityRepo) GetAll(ctx context.Context) ([]*entity.Community, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*entity.Community, 0, len(m.communities))
	for _, comm := range m.communities {
		result = append(result, comm)
	}
	return result, nil
}

// Delete: コミュニティを削除
func (m *MemoryCommunityRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.communities[id]; !ok {
		return errors.New("community not found")
	}
	delete(m.communities, id)
	return nil
}

// インタフェース実装をチェック
var _ repository.CommunityRepository = (*MemoryCommunityRepo)(nil)
