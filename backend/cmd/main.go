package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/infrastructure/repository"
	"github.com/rayfiyo/zousui/backend/interface/controller"
	"github.com/rayfiyo/zousui/backend/interface/gateway"
	"github.com/rayfiyo/zousui/backend/interface/router"
	"github.com/rayfiyo/zousui/backend/usecase"
	"github.com/rayfiyo/zousui/backend/utils/config"
)

func main() {
	// リポジトリ初期化
	communityRepo := repository.NewMemoryCommunityRepo()
	agentRepo := repository.NewMemoryAgentRepo()

	// 環境変数ロード
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	// Gemini ゲートウェイ
	llmGw, err := gateway.NewGeminiLLMGateway(context.Background())
	if err != nil {
		log.Fatalf("failed to create gemini gateway: %v", err)
	}
	defer llmGw.Client.Close()

	// Mockゲートウェイ
	mockGw := &gateway.MockLLMGatewayJSON{}

	// 既存: シミュレーション/外交/コミュニティユースケース
	simulateUC := usecase.NewSimulateCultureEvolutionUsecase(communityRepo, agentRepo, llmGw)
	diploUC := usecase.NewDiplomacyUsecase(communityRepo, llmGw)
	communityUC := usecase.NewCommunityUsecase(communityRepo)

	// 集約ゲートウェイ (複数LLMを内部でランダム使用するサンプル)
	multiGw := gateway.NewMultiLLMGateway(llmGw, mockGw)

	// 新規: コミュニティ同士の干渉ユースケース
	interferenceBetweenCommunitiesUC := usecase.NewSimulateInterferenceBetweenCommunitiesUsecase(
		communityRepo,
		multiGw, // 複数LLMを渡す
	)

	// コントローラ
	commCtrl := controller.NewCommunityController(communityUC)
	diploCtrl := controller.NewDiplomacyController(diploUC)
	simCtrl := controller.NewSimulateController(simulateUC)
	imageCtrl := controller.NewImageController(*communityUC)

	// 干渉コントローラ
	interferenceCtrl := controller.NewInterferenceController(interferenceBetweenCommunitiesUC)

	// データ初期化
	seedData(communityRepo, agentRepo)

	// ルーティング
	r := router.NewRouter(
		commCtrl,
		diploCtrl,
		simCtrl,
		imageCtrl,
		interferenceCtrl, // 新たに渡す
	)

	fmt.Println("Starting zousui MVP server on :8080")
	r.Run(":8080")
}

// seedData: テスト用の初期データを挿入
func seedData(cr *repository.MemoryCommunityRepo, ar *repository.MemoryAgentRepo) {
	communityID := "comm-1"
	comm := &entity.Community{
		ID:         communityID,
		Name:       "DesertTribe",
		Population: 100,
		Culture:    "砂漠での生存術が中心の文化",
	}
	cr.Save(context.TODO(), comm)

	// 他に複数作ってもOK
	communityID2 := "comm-2"
	comm2 := &entity.Community{
		ID:         communityID2,
		Name:       "OceanicCity",
		Population: 300,
		Culture:    "海底で歌と踊りを好む平和な国",
	}
	cr.Save(context.TODO(), comm2)

	agent1 := &entity.Agent{
		ID:          "agent-1",
		Name:        "Aisha",
		CommunityID: communityID,
		Personality: "好奇心旺盛で穏やか",
	}
	agent2 := &entity.Agent{
		ID:          "agent-2",
		Name:        "Jamal",
		CommunityID: communityID,
		Personality: "勇敢で戦闘的",
	}

	ar.Agents = append(ar.Agents, agent1, agent2)
}
