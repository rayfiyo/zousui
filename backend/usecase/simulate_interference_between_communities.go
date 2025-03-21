package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"github.com/rayfiyo/zousui/backend/utils/consts"
	"go.uber.org/zap"
)

type SimulateInterferenceBetweenCommunitiesUsecase struct {
	communityRepo  repository.CommunityRepository
	llmGateway     repository.LLMGateway
	simulationRepo repository.SimulationRepository
}

func NewSimulateInterferenceBetweenCommunitiesUsecase(
	cr repository.CommunityRepository,
	lg repository.LLMGateway,
	sr repository.SimulationRepository,
) *SimulateInterferenceBetweenCommunitiesUsecase {
	zap.L().Debug("Initializing SimulateInterferenceBetweenCommunitiesUsecase")
	return &SimulateInterferenceBetweenCommunitiesUsecase{
		communityRepo:  cr,
		llmGateway:     lg,
		simulationRepo: sr,
	}
}

// 干渉シミュレーション本体
func (uc *SimulateInterferenceBetweenCommunitiesUsecase) Execute(
	ctx context.Context,
	commAID, commBID, userInput string,
) error {
	logger := zap.L()

	// コミュニティA,Bを取得
	logger.Debug("Starting interference between communities",
		zap.String("commA", commAID), zap.String("commB", commBID))
	commA, err := uc.communityRepo.GetByID(ctx, commAID)
	if err != nil {
		logger.Error("Failed to get community A",
			zap.String("commA", commAID), zap.Error(err))
		return fmt.Errorf("failed to get community A: %w", err)
	}
	commB, err := uc.communityRepo.GetByID(ctx, commBID)
	if err != nil {
		logger.Error("Failed to get community B",
			zap.String("commB", commBID), zap.Error(err))
		return fmt.Errorf("failed to get community B: %w", err)
	}

	// プロンプト作成: 2つのコミュニティの文化が互いに干渉したらどうなるか
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

	// LLM呼び出し (MultiLLMGatewayを想定)
	logger.Debug("Interference between communities prompt",
		zap.String("prompt", prompt))
	llmResp, err := uc.llmGateway.GenerateCultureUpdate(ctx, prompt, userInput)
	if err != nil {
		logger.Error("LLM generation failed", zap.Error(err))
		return fmt.Errorf("failed to generate culture update: %w", err)
	}

	// JSONパース
	var result struct {
		NewCultureA       string `json:"newCultureA"`
		PopulationChangeA int    `json:"populationChangeA"`
		NewCultureB       string `json:"newCultureB"`
		PopulationChangeB int    `json:"populationChangeB"`
	}
	if err := json.Unmarshal([]byte(llmResp), &result); err != nil {
		logger.Error("JSON unmarshal failed", zap.Error(err),
			zap.String("response", llmResp))
		return fmt.Errorf("invalid JSON from LLM: %w\nLLM response was: %s",
			err, llmResp)
	}

	// expected フィールドが空なら、"newCulture" と "populationChange" キーで再取得
	if result.NewCultureA == "" && result.NewCultureB == "" {
		logger.Warn("LLM response missing expected keys, applying fallback parsing")
		var tmp map[string]interface{}
		if err2 := json.Unmarshal([]byte(llmResp), &tmp); err2 == nil {
			if nc, ok := tmp["newCulture"].(string); ok {
				// A のみ変更を適用
				result.NewCultureA = nc
				// result.NewCultureB = nc
			}
			if pc, ok := tmp["populationChange"].(float64); ok {
				result.PopulationChangeA = int(pc)
				// result.PopulationChangeB = int(pc)
			}
		}
	}
	logger.Debug("Interference result", zap.Any("result", result))

	// 結果をコミュニティA, Bに反映
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

	//
	if err := uc.communityRepo.Save(ctx, commA); err != nil {
		logger.Error("Failed to save community A",
			zap.String("commA", commAID), zap.Error(err))
		return fmt.Errorf("failed to save community A: %w", err)
	}
	if err := uc.communityRepo.Save(ctx, commB); err != nil {
		logger.Error("Failed to save community B",
			zap.String("commB", commBID), zap.Error(err))
		return fmt.Errorf("failed to save community B: %w", err)
	}

	resultJSON, _ := json.Marshal(result)
	simResult := &entity.SimulationResult{
		Type:        "interference",
		Communities: []string{commAID, commBID},
		ResultJSON:  string(resultJSON),
	}
	if err := uc.simulationRepo.Save(ctx, simResult); err != nil {
		logger.Error("Failed to save community B",
			zap.String("commB", commBID), zap.Error(err))
		return fmt.Errorf("failed to save community B: %w", err)
	}

	logger.Info("Interference between communities executed successfully",
		zap.String("commA", commAID), zap.String("commB", commBID))
	return nil
}
