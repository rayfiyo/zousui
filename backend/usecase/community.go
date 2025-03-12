package usecase

import (
	"context"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type CommunityUsecase struct {
	communityRepo repository.CommunityRepository
}

func NewCommunityUsecase(
	repo repository.CommunityRepository,
) *CommunityUsecase {
	zap.L().Debug("Initializing CommunityUsecase")
	return &CommunityUsecase{
		communityRepo: repo,
	}
}

// コミュニティを作成
func (cu *CommunityUsecase) CreateCommunity(
	ctx context.Context,
	comm *entity.Community,
) error {
	logger := zap.L()

	logger.Debug("Creating community", zap.String("communityID", comm.ID))

	// IDが重複していないか簡単にチェック
	existing, _ := cu.communityRepo.GetByID(ctx, comm.ID)
	if existing != nil {
		logger.Warn("Community already exists", zap.String("communityID", comm.ID))
		return fmt.Errorf("community ID %s already exists", comm.ID)
	}

	// Save
	if err := cu.communityRepo.Save(ctx, comm); err != nil {
		logger.Error("Failed to save community",
			zap.String("communityID", comm.ID), zap.Error(err))
		return err
	}

	logger.Info("Community created", zap.String("communityID", comm.ID))
	return nil
}

// コミュニティを更新
func (cu *CommunityUsecase) GetCommunityByID(
	ctx context.Context,
	id string,
) (*entity.Community, error) {
	zap.L().Debug("Fetching community by ID", zap.String("communityID", id))
	return cu.communityRepo.GetByID(ctx, id)
}

// コミュニティを更新
func (cu *CommunityUsecase) DeleteCommunity(ctx context.Context, id string) error {
	logger := zap.L()
	logger.Debug("Deleting community", zap.String("communityID", id))
	err := cu.communityRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("Failed to delete community", zap.String("communityID", id), zap.Error(err))
	} else {
		logger.Info("Community deleted", zap.String("communityID", id))
	}
	return err
}

// 全コミュニティを取得
func (cu *CommunityUsecase) GetAllCommunities(
	ctx context.Context,
) ([]*entity.Community, error) {
	zap.L().Debug("Getting all communities")
	return cu.communityRepo.GetAll(ctx)
}
