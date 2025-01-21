package usecase

import (
	"context"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
)

// CommunityUsecase: コミュニティに関するユースケースをまとめる
type CommunityUsecase struct {
	communityRepo repository.CommunityRepository
}

func NewCommunityUsecase(repo repository.CommunityRepository) *CommunityUsecase {
	return &CommunityUsecase{
		communityRepo: repo,
	}
}

// CreateCommunity: コミュニティを作成
func (cu *CommunityUsecase) CreateCommunity(ctx context.Context, comm *domain.Community) error {
	// IDが重複していないか簡単にチェック
	existing, _ := cu.communityRepo.GetByID(ctx, comm.ID)
	if existing != nil {
		return fmt.Errorf("community ID %s already exists", comm.ID)
	}
	// Save
	if err := cu.communityRepo.Save(ctx, comm); err != nil {
		return err
	}
	return nil
}

// UpdateCommunity: コミュニティを更新
func (cu *CommunityUsecase) GetCommunityByID(ctx context.Context, id string) (*domain.Community, error) {
	return cu.communityRepo.GetByID(ctx, id)
}

// UpdateCommunity: コミュニティを更新
func (cu *CommunityUsecase) DeleteCommunity(ctx context.Context, id string) error {
	return cu.communityRepo.Delete(ctx, id)
}

// GetAllCommunities: 全コミュニティを取得
func (cu *CommunityUsecase) GetAllCommunities(ctx context.Context) ([]*entity.Community, error) {
	// リポジトリに全件取得メソッドを用意する or InMemoryなら鍵をすべて走査
	return cu.communityRepo.GetAll(ctx)
}
