package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/repository"
)

type DiplomacyUsecase struct {
	communityRepo repository.CommunityRepository
	llmGateway    repository.LLMGateway
}

func NewDiplomacyUsecase(cr repository.CommunityRepository, lg repository.LLMGateway) *DiplomacyUsecase {
	return &DiplomacyUsecase{communityRepo: cr, llmGateway: lg}
}

func (du *DiplomacyUsecase) ExecuteDiplomacy(ctx context.Context, commAID, commBID string) error {
	commA, err := du.communityRepo.GetByID(ctx, commAID)
	if err != nil {
		return fmt.Errorf("failed to get community A: %w", err)
	}
	commB, err := du.communityRepo.GetByID(ctx, commBID)
	if err != nil {
		return fmt.Errorf("failed to get community B: %w", err)
	}

	// プロンプト
	prompt := fmt.Sprintf(`コミュニティA: {Name: %s, Population: %d, Culture: %s}
コミュニティB: {Name: %s, Population: %d, Culture: %s}
この2つのコミュニティが外交交渉を行い、その結果をJSONで返してください。
必ず以下の形式に従うこと:

{
  "outcome": "peace|war|trade|alliance",
  "description": "交渉の結果・内容",
  "popChangeA": 0,
  "popChangeB": 0
}
`,
		commA.Name, commA.Population, commA.Culture,
		commB.Name, commB.Population, commB.Culture)

	llmResp, err := du.llmGateway.GenerateCultureUpdate(ctx, prompt)
	if err != nil {
		return err
	}

	// JSONパース
	var result struct {
		Outcome     string `json:"outcome"`
		Description string `json:"description"`
		PopChangeA  int    `json:"popChangeA"`
		PopChangeB  int    `json:"popChangeB"`
	}
	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
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

	// 文化欄にも何か反映する例
	commA.UpdateCulture(fmt.Sprintf("%s | %s", commA.Culture, result.Outcome))
	commB.UpdateCulture(fmt.Sprintf("%s | %s", commB.Culture, result.Outcome))

	// 保存
	if err := du.communityRepo.Save(ctx, commA); err != nil {
		return err
	}
	if err := du.communityRepo.Save(ctx, commB); err != nil {
		return err
	}

	return nil
}
