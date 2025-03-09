package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type SimulateCultureEvolutionUsecase struct {
	communityRepo repository.CommunityRepository
	agentRepo     repository.AgentRepository
	llmGateway    repository.LLMGateway
}

func NewSimulateCultureEvolutionUsecase(
	cr repository.CommunityRepository,
	ar repository.AgentRepository,
	lg repository.LLMGateway,
) *SimulateCultureEvolutionUsecase {
	zap.L().Debug("Initializing SimulateCultureEvolutionUsecase")
	return &SimulateCultureEvolutionUsecase{
		communityRepo: cr,
		agentRepo:     ar,
		llmGateway:    lg,
	}
}

// コミュニティを指定して、エージェントとLLMを用いた文化進化シミュレーションを実行する
func (uc *SimulateCultureEvolutionUsecase) Execute(
	ctx context.Context,
	communityID string,
) error {
	logger := zap.L()

	// コミュニティを取得
	logger.Debug("Starting culture evolution simulation",
		zap.String("communityID", communityID))
	comm, err := uc.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		logger.Error("Failed to get community",
			zap.String("communityID", communityID), zap.Error(err))
		return fmt.Errorf("failed to get community: %w", err)
	}

	// コミュニティに所属するエージェントを取得
	agents, err := uc.agentRepo.GetAgentsByCommunity(ctx, communityID)
	if err != nil {
		logger.Error("Failed to get agents",
			zap.String("communityID", communityID), zap.Error(err))
		return fmt.Errorf("failed to get agents: %w", err)
	}

	// シミュレーション用のプロンプトを作成
	prompt := fmt.Sprintf(
		"コミュニティ名: {{%s}}\n人口: {{%d}}\n現文化: {{%s}}\n---\n",
		comm.Name, comm.Population, comm.Culture,
	)
	for _, agent := range agents {
		prompt += fmt.Sprintf(
			"エージェント: %s, 性格: %s\n", agent.Name, agent.Personality,
		)
	}

	// LLMに問い合わせ
	logger.Debug("Simulation prompt", zap.String("prompt", prompt))
	llmResp, err := uc.llmGateway.GenerateCultureUpdate(ctx, prompt, "")
	if err != nil {
		logger.Error("LLM generation failed", zap.Error(err))
		return fmt.Errorf("failed to generate culture update: %w", err)
	}

	// JSONパース
	logger.Debug("LLM response", zap.String("response", llmResp))
	var result entity.CultureUpdateResponse
	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
		logger.Error("JSON unmarshal failed", zap.Error(err))
		return fmt.Errorf("invalid JSON from LLM (failed to parse LLM JSON): %w", err)
	}

	// ドメインモデルを使って更新
	comm.UpdateCulture(result.NewCulture)
	comm.Population += result.PopulationChange
	if err := uc.communityRepo.Save(ctx, comm); err != nil {
		logger.Error("Failed to save community after simulation",
			zap.String("communityID", communityID), zap.Error(err))
		return fmt.Errorf("failed to save community: %w", err)
	}
	logger.Info("Culture evolution simulation executed successfully",
		zap.String("communityID", communityID))

	return nil
}
