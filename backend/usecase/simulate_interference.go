package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type SimulateInterferenceUsecase struct {
	communityRepo repository.CommunityRepository
	agentRepo     repository.AgentRepository
	llmGateway    repository.LLMGateway
}

func NewSimulateInterferenceUsecase(
	cr repository.CommunityRepository,
	ar repository.AgentRepository,
	lg repository.LLMGateway,
) *SimulateInterferenceUsecase {
	zap.L().Debug("Initializing SimulateInterferenceUsecase")
	return &SimulateInterferenceUsecase{
		communityRepo: cr,
		agentRepo:     ar,
		llmGateway:    lg,
	}
}

// 干渉シナリオを実行する
func (uc *SimulateInterferenceUsecase) Execute(
	ctx context.Context,
	communityID string,
) error {
	logger := zap.L()

	// コミュニティを取得
	logger.Debug("Starting interference simulation",
		zap.String("communityID", communityID))
	comm, err := uc.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		logger.Error("Failed to get community",
			zap.String("communityID", communityID), zap.Error(err))
		return fmt.Errorf("failed to get community: %w", err)
	}

	// 関連エージェントを取得
	agents, err := uc.agentRepo.GetAgentsByCommunity(ctx, communityID)
	if err != nil {
		logger.Error("Failed to get agents",
			zap.String("communityID", communityID), zap.Error(err))
		return fmt.Errorf("failed to get agents: %w", err)
	}

	// 干渉シナリオ用のプロンプト作成例
	prompt := fmt.Sprintf(
		"コミュニティ名: %q\n人口: %d\n文化: %q\n---\nエージェント情報:\n",
		comm.Name, comm.Population, comm.Culture,
	)
	for _, a := range agents {
		prompt += fmt.Sprintf("- %s (性格: %s)\n", a.Name, a.Personality)
	}

	prompt += `ここに対して、複数の知性(LLM)から干渉アイデアが持ち込まれました。
    新しい文化アイデアや予想外の変化を考えて、必ず以下のJSON形式を返してください:
    {
      "newCulture": "...",
      "populationChange": 0
    }`

	// LLM呼び出し (MultiLLMGateway などが内部で複数LLMを利用)
	logger.Debug("Interference simulation prompt", zap.String("prompt", prompt))
	llmResp, err := uc.llmGateway.GenerateCultureUpdate(ctx, prompt, "")
	if err != nil {
		return fmt.Errorf("failed to generate culture update: %w", err)
	}

	// JSONパース
	logger.Debug("LLM interference response", zap.String("response", llmResp))
	var result entity.CultureUpdateResponse
	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
		logger.Warn("JSON unmarshal failed, using raw response", zap.Error(err))
		comm.UpdateCulture(comm.Culture + " | " + llmResp)
	} else {
		comm.UpdateCulture(result.NewCulture)
		comm.Population += result.PopulationChange
	}

	// 人口が0未満にならないように
	if comm.Population < 0 {
		comm.Population = 0
	}

	// DBへ保存
	if err := uc.communityRepo.Save(ctx, comm); err != nil {
		logger.Error("Failed to save updated community",
			zap.String("communityID", communityID), zap.Error(err))
		return fmt.Errorf("failed to save updated community: %w", err)
	}

	logger.Info("Interference simulation executed successfully",
		zap.String("communityID", communityID))
	return nil
}
