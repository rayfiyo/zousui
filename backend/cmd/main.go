package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/domain"
	"github.com/rayfiyo/zousui/backend/infrastructure/repository"
	"github.com/rayfiyo/zousui/backend/interface/controller"
	"github.com/rayfiyo/zousui/backend/interface/gateway"
	"github.com/rayfiyo/zousui/backend/usecase"
)

func main() {
	// =====================================
	// 依存オブジェクト(リポジトリ & ゲートウェイ)初期化
	// =====================================
	communityRepo := repository.NewMemoryCommunityRepo()
	agentRepo := repository.NewMemoryAgentRepo()

	// 例として モックLLMゲートウェイ
	llmGw := &gateway.MockLLMGateway{}

	// =====================================
	// ユースケース生成
	// =====================================
	simulateUC := usecase.NewSimulateCultureEvolutionUsecase(communityRepo, agentRepo, llmGw)

	// =====================================
	// データ初期投入 (コミュニティやエージェント)
	// =====================================
	seedData(communityRepo, agentRepo)

	// =====================================
	// HTTPサーバ起動
	// =====================================
	r := gin.Default()

	// コントローラ生成 & ルート設定
	simCtrl := controller.NewSimulateController(simulateUC)
	simCtrl.SetupRoutes(r)

	fmt.Println("Starting zousui MVP server on :8080")
	r.Run(":8080")
}

// seedData: テスト用の初期データを挿入
func seedData(cr *repository.MemoryCommunityRepo, ar *repository.MemoryAgentRepo) {
	communityID := "comm-1"
	comm := &domain.Community{
		ID:         communityID,
		Name:       "DesertTribe",
		Population: 100,
		Culture:    "砂漠での生存術が中心の文化",
	}
	cr.Save(nil, comm)

	agent1 := &domain.Agent{
		ID:          "agent-1",
		Name:        "Aisha",
		CommunityID: communityID,
		Personality: "好奇心旺盛で穏やか",
	}
	agent2 := &domain.Agent{
		ID:          "agent-2",
		Name:        "Jamal",
		CommunityID: communityID,
		Personality: "勇敢で戦闘的",
	}

	ar.Agents = append(ar.Agents, agent1, agent2)
}
