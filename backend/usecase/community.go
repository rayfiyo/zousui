package usecase

import (
	"context"

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

// GetAllCommunities: 全コミュニティを取得
func (cu *CommunityUsecase) GetAllCommunities(ctx context.Context) ([]*entity.Community, error) {
	// リポジトリに全件取得メソッドを用意する or InMemoryなら鍵をすべて走査
	return cu.communityRepo.GetAll(ctx)
}
