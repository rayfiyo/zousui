package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
)

// SimulateInterferenceUsecase:
// 複数LLMの結果を「干渉」として取り込み、コミュニティに反映させるシミュレーション例
type SimulateInterferenceUsecase struct {
	communityRepo repository.CommunityRepository
	agentRepo     repository.AgentRepository
	llmGateway    repository.LLMGateway
}

// NewSimulateInterferenceUsecase: コンストラクタ
func NewSimulateInterferenceUsecase(
	cr repository.CommunityRepository,
	ar repository.AgentRepository,
	lg repository.LLMGateway,
) *SimulateInterferenceUsecase {
	return &SimulateInterferenceUsecase{
		communityRepo: cr,
		agentRepo:     ar,
		llmGateway:    lg,
	}
}

// Execute: 干渉シミュレーション
//   - コミュニティとエージェントの情報をまとめてプロンプト
//   - 集約ゲートウェイ(MultiLLMGateway)などを経由して複数LLMから回答
//   - 結果(JSON)をパースし、コミュニティへ反映
func (uc *SimulateInterferenceUsecase) Execute(ctx context.Context, communityID string) error {
	// コミュニティを取得
	comm, err := uc.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return fmt.Errorf("failed to get community: %w", err)
	}

	// 関連エージェントを取得
	agents, err := uc.agentRepo.GetAgentsByCommunity(ctx, communityID)
	if err != nil {
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

	prompt += `
    ここに対して、複数の知性(LLM)から干渉アイデアが持ち込まれました。
    新しい文化アイデアや予想外の変化を考えて、必ず以下のJSON形式を返してください:
    {
      "newCulture": "...",
      "populationChange": 0
    }
    `

	// LLM呼び出し (MultiLLMGateway などが内部で複数LLMを利用)
	llmResp, err := uc.llmGateway.GenerateCultureUpdate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate culture update: %w", err)
	}

	// JSONパース
	var result entity.CultureUpdateResponse
	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
		// JSONパースに失敗した場合は、そのまま文字列を文化にしてしまう例
		// 干渉で壊れたJSONを受容するイメージ
		comm.UpdateCulture(comm.Culture + " | " + llmResp)
	} else {
		// 正常にパースできた場合は通常処理
		comm.UpdateCulture(result.NewCulture)
		comm.Population += result.PopulationChange
	}

	// 人口が0未満にならないように
	if comm.Population < 0 {
		comm.Population = 0
	}

	// DBへ保存
	if err := uc.communityRepo.Save(ctx, comm); err != nil {
		return fmt.Errorf("failed to save updated community: %w", err)
	}

	return nil
}
