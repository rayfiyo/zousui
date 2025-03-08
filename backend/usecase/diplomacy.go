package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type DiplomacyUsecase struct {
	communityRepo repository.CommunityRepository
	llmGateway    repository.LLMGateway
}

func NewDiplomacyUsecase(
	cr repository.CommunityRepository,
	lg repository.LLMGateway,
) *DiplomacyUsecase {
	zap.L().Debug("Initializing DiplomacyUsecase")
	return &DiplomacyUsecase{communityRepo: cr, llmGateway: lg}
}

// 2つのコミュニティ間の外交交渉を実行する
func (du *DiplomacyUsecase) ExecuteDiplomacy(
	ctx context.Context,
	commAID, commBID string,
) error {
	logger := zap.L()

    // コミュニティを取得
	logger.Debug("Executing diplomacy simulation",
		zap.String("commA", commAID), zap.String("commB", commBID))
	commA, err := du.communityRepo.GetByID(ctx, commAID)
	if err != nil {
		return fmt.Errorf("failed to get community A: %w", err)
	}
	commB, err := du.communityRepo.GetByID(ctx, commBID)
	if err != nil {
		return fmt.Errorf("failed to get community B: %w", err)
	}

    // LLMに外交交渉をリクエスト
	prompt := fmt.Sprintf(`コミュニティA: {Name: %s, Population: %d, Culture: %s}
    コミュニティB: {Name: %s, Population: %d, Culture: %s}
    この2つのコミュニティが外交交渉を行い、その結果をJSONで返してください。
    必ず以下の形式に従うこと:
    {
      "outcome": "peace|war|trade|alliance",
      "description": "交渉の結果・内容",
      "popChangeA": 0,
      "popChangeB": 0
    }`,
		commA.Name, commA.Population, commA.Culture,
		commB.Name, commB.Population, commB.Culture)

    // LLMにリクエスト
	logger.Debug("Diplomacy prompt", zap.String("prompt", prompt))
	llmResp, err := du.llmGateway.GenerateCultureUpdate(ctx, prompt)
	if err != nil {
		logger.Error("LLM generation failed", zap.Error(err))
		return err
	}

	// JSONパース
	logger.Debug("LLM response received", zap.String("response", llmResp))
	var result struct {
		Outcome     string `json:"outcome"`
		Description string `json:"description"`
		PopChangeA  int    `json:"popChangeA"`
		PopChangeB  int    `json:"popChangeB"`
	}
	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
		logger.Error("JSON unmarshal failed", zap.Error(err))
		return fmt.Errorf("invalid JSON from LLM: %w", err)
	}

	// コミュニティに反映
	commA.Population += result.PopChangeA
	if commA.Population < 0 {
		commA.Population = 0
	}
	commB.Population += result.PopChangeB
	if commB.Population < 0 {
		commB.Population = 0
	}

    // 結果を記録
	commA.UpdateCulture(fmt.Sprintf("%s | %s", commA.Culture, result.Outcome))
	commB.UpdateCulture(fmt.Sprintf("%s | %s", commB.Culture, result.Outcome))

	// 保存
	if err := du.communityRepo.Save(ctx, commA); err != nil {
		logger.Error("Failed to save community A",
			zap.String("commA", commAID), zap.Error(err))
		return err
	}
	if err := du.communityRepo.Save(ctx, commB); err != nil {
		logger.Error("Failed to save community B",
			zap.String("commB", commBID), zap.Error(err))
		return err
	}
	logger.Info("Diplomacy simulation executed successfully",
		zap.String("commA", commAID), zap.String("commB", commBID))
	return nil
}
