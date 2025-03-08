package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/repository"
	"github.com/rayfiyo/zousui/backend/utils/consts"
)

// 2つのコミュニティ同士が干渉(相互作用)するシミュレーションを行う。
// 複数LLMのゲートウェイを呼び出して、その応答を反映させる。
type SimulateInterferenceBetweenCommunitiesUsecase struct {
	communityRepo repository.CommunityRepository
	llmGateway    repository.LLMGateway // MultiLLMGateway を想定
}

// コンストラクタ
func NewSimulateInterferenceBetweenCommunitiesUsecase(
	cr repository.CommunityRepository,
	lg repository.LLMGateway,
) *SimulateInterferenceBetweenCommunitiesUsecase {
	return &SimulateInterferenceBetweenCommunitiesUsecase{
		communityRepo: cr,
		llmGateway:    lg,
	}
}

// 干渉シミュレーション本体
//   - commA, commB: 干渉させる2つのコミュニティID
func (uc *SimulateInterferenceBetweenCommunitiesUsecase) Execute(
	ctx context.Context,
	commAID, commBID string,
) error {
	// 1. コミュニティA,Bを取得
	commA, err := uc.communityRepo.GetByID(ctx, commAID)
	if err != nil {
		return fmt.Errorf("failed to get community A: %w", err)
	}
	commB, err := uc.communityRepo.GetByID(ctx, commBID)
	if err != nil {
		return fmt.Errorf("failed to get community B: %w", err)
	}

	// 2. プロンプト作成: 2つのコミュニティの文化が互いに干渉したらどうなるか
	prompt := fmt.Sprintf(`
        これは2つのコミュニティA,Bが「干渉」しあうシミュレーションです。
        コミュニティA:
          Name: %s
          Population: %d
          Culture: %s

        コミュニティB:
          Name: %s
          Population: %d
          Culture: %s

        これらが互いの文化に影響を与え合った結果、どのように変化するかを予測し、
        必ず下記JSON形式で出力してください:
        {
          "newCultureA": "string",
          "populationChangeA": 0,
          "newCultureB": "string",
          "populationChangeB": 0
        }

        出力例:
        {
          "newCultureA": "AがBから得た刺激を表す新文化",
          "populationChangeA": 5,
          "newCultureB": "BがAから得た衝撃を表す新文化",
          "populationChangeB": -2
        }

        %s`,
		commA.Name, commA.Population, commA.Culture,
		commB.Name, commB.Population, commB.Culture,
		consts.SpecifyingResponseFormat, // 追加で「必ずJSONで返してね」の指示
	)

	// 3. LLM呼び出し (MultiLLMGatewayを想定)
	llmResp, err := uc.llmGateway.GenerateCultureUpdate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate culture update: %w", err)
	}

	// 4. JSONパース
	var result struct {
		NewCultureA       string `json:"newCultureA"`
		PopulationChangeA int    `json:"populationChangeA"`
		NewCultureB       string `json:"newCultureB"`
		PopulationChangeB int    `json:"populationChangeB"`
	}

	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
		// JSONパース失敗時は強制的にエラー
		return fmt.Errorf("invalid JSON from LLM: %w\nLLM response was: %s",
			err, llmResp)
	}

	// 5. 結果をコミュニティA, Bに反映
	//    文化と人口を変更
	if result.NewCultureA != "" {
		commA.Culture = result.NewCultureA
	}
	commA.Population += result.PopulationChangeA
	if commA.Population < 0 {
		commA.Population = 0
	}

	if result.NewCultureB != "" {
		commB.Culture = result.NewCultureB
	}
	commB.Population += result.PopulationChangeB
	if commB.Population < 0 {
		commB.Population = 0
	}

	// 6. 保存
	if err := uc.communityRepo.Save(ctx, commA); err != nil {
		return fmt.Errorf("failed to save community A: %w", err)
	}
	if err := uc.communityRepo.Save(ctx, commB); err != nil {
		return fmt.Errorf("failed to save community B: %w", err)
	}

	return nil
}
