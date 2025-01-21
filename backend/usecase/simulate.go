package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain"
	"github.com/rayfiyo/zousui/backend/utils"
)

// AgentRepository: エージェントに関するリポジトリインタフェース(読み書き)
type AgentRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Agent, error)
	Save(ctx context.Context, agent *domain.Agent) error
	GetAll(ctx context.Context) ([]*domain.Agent, error)
	GetAgentsByCommunity(ctx context.Context, communityID string) ([]*domain.Agent, error)
}

// CommunityRepository: コミュニティに関するリポジトリインタフェース(読み書き)
type CommunityRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Community, error)
	Save(ctx context.Context, community *domain.Community) error
	GetAll(ctx context.Context) ([]*domain.Community, error)
}

// LLMGateway: LLMに問い合わせるためのインタフェース
type LLMGateway interface {
	GenerateCultureUpdate(ctx context.Context, prompt string) (string, error)
}

// SimulateCultureEvolutionUsecase: 文化を進化(変化)させるシミュレーション例
type SimulateCultureEvolutionUsecase struct {
	communityRepo CommunityRepository
	agentRepo     AgentRepository
	llmGateway    LLMGateway
}

// NewSimulateCultureEvolutionUsecase: コンストラクタ
func NewSimulateCultureEvolutionUsecase(
	cr CommunityRepository,
	ar AgentRepository,
	lg LLMGateway,
) *SimulateCultureEvolutionUsecase {
	return &SimulateCultureEvolutionUsecase{
		communityRepo: cr,
		agentRepo:     ar,
		llmGateway:    lg,
	}
}

// Execute: コミュニティを指定して、エージェントとLLMを用いた文化進化シミュレーションを実行する
func (uc *SimulateCultureEvolutionUsecase) Execute(ctx context.Context, communityID string) error {
	// コミュニティを取得
	comm, err := uc.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return fmt.Errorf("failed to get community: %w", err)
	}

	// 関連するエージェントを取得
	agents, err := uc.agentRepo.GetAgentsByCommunity(ctx, communityID)
	if err != nil {
		return fmt.Errorf("failed to get agents: %w", err)
	}

	// 簡単な例: エージェントの情報からプロンプトを組み立てる
	prompt := fmt.Sprintf(
		"コミュニティ名: {{%s}}\n人口: {{%d}}\n現文化: {{%s}}\n%s\n---\n",
		comm.Name, comm.Population, comm.Culture, utils.SpecifyingResponseFormat,
	)
	for _, agent := range agents {
		prompt += fmt.Sprintf(
			"エージェント: %s, 性格: %s\n", agent.Name, agent.Personality,
		)
	}

	// LLMに問い合わせ
	llmResp, err := uc.llmGateway.GenerateCultureUpdate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate culture update: %w", err)
	}

	// JSONパース
	var result domain.CultureUpdateResponse
	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
		// ここで 単純に文章全体を newCulture に入れる処理もアリ
		// result.NewCulture = llmResp
		// result.PopulationChange = 0
		return fmt.Errorf("invalid JSON from LLM (failed to parse LLM JSON): %w", err)
	}

	// ドメインモデルを使って更新
	comm.UpdateCulture(result.NewCulture)
	comm.Population += result.PopulationChange

	// コミュニティを保存
	if err := uc.communityRepo.Save(ctx, comm); err != nil {
		return fmt.Errorf("failed to save community: %w", err)
	}

	return nil
}
