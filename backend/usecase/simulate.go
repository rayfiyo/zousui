package usecase

import (
	"context"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain"
)

// CommunityRepository: コミュニティに関するリポジトリインタフェース(読み書き)
type CommunityRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Community, error)
	Save(ctx context.Context, community *domain.Community) error
	GetAll(ctx context.Context) ([]*domain.Community, error)
}

// AgentRepository: エージェントに関するリポジトリインタフェース(読み書き)
type AgentRepository interface {
	GetAgentsByCommunity(ctx context.Context, communityID string) ([]*domain.Agent, error)
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
	prompt := fmt.Sprintf("コミュニティ名: %s\n人口: %d\n現文化: %s\n---\n",
		comm.Name, comm.Population, comm.Culture)

	for _, agent := range agents {
		prompt += fmt.Sprintf("エージェント: %s, 性格: %s\n", agent.Name, agent.Personality)
	}
	prompt += "このコミュニティの文化を新しい方向に進化させるアイデアを提案してください。"

	// LLMに問い合わせ
	newCulture, err := uc.llmGateway.GenerateCultureUpdate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate culture update: %w", err)
	}

	// ドメインモデルを使って更新
	comm.UpdateCulture(newCulture)

	// コミュニティを保存
	if err := uc.communityRepo.Save(ctx, comm); err != nil {
		return fmt.Errorf("failed to save community: %w", err)
	}

	return nil
}
