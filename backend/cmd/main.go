package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/infrastructure/repository"
	"github.com/rayfiyo/zousui/backend/interface/controller"
	"github.com/rayfiyo/zousui/backend/interface/gateway"
	"github.com/rayfiyo/zousui/backend/interface/router"
	"github.com/rayfiyo/zousui/backend/usecase"
	"github.com/rayfiyo/zousui/backend/utils/config"
	"go.uber.org/zap"
)

func main() {
	// zap ロガーの初期化
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // ログバッファのフラッシュ
	zap.ReplaceGlobals(logger)

	// リポジトリ初期化
	communityRepo := repository.NewMemoryCommunityRepo()
	agentRepo := repository.NewMemoryAgentRepo()
	logger.Debug("Repositories initialized")

	// 環境変数ロード
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	logger.Info("Environment variables loaded")

	// Gemini ゲートウェイ
	llmGw, err := gateway.NewGeminiLLMGateway(context.Background())
	if err != nil {
		logger.Fatal("failed to create gemini gateway", zap.Error(err))
	}
	defer llmGw.Client.Close()
	logger.Info("Gemini gateway created")

	// Mockゲートウェイ
	mockGw := &gateway.MockLLMGatewayJSON{}
	logger.Debug("Mock gateway initialized")

	// シミュレーション/外交/コミュニティユースケース
	simulateUC := usecase.NewSimulateCultureEvolutionUsecase(communityRepo, agentRepo, llmGw)
	diploUC := usecase.NewDiplomacyUsecase(communityRepo, llmGw)
	communityUC := usecase.NewCommunityUsecase(communityRepo)
	logger.Debug("Usecases initialized")

	// 集約ゲートウェイ (複数LLMを内部でランダム使用するサンプル)
	multiGw := gateway.NewMultiLLMGateway(llmGw, mockGw)
	logger.Debug("Multi LLM gateway initialized")

	// コミュニティ同士の干渉ユースケース
	interferenceBetweenCommunitiesUC := usecase.
		NewSimulateInterferenceBetweenCommunitiesUsecase(communityRepo, multiGw)
	logger.Debug("Interference usecase initialized")

	// コントローラ
	commCtrl := controller.NewCommunityController(communityUC)
	diploCtrl := controller.NewDiplomacyController(diploUC)
	simCtrl := controller.NewSimulateController(simulateUC)
	imageCtrl := controller.NewImageController(*communityUC)
	interferenceCtrl := controller.NewInterferenceController(interferenceBetweenCommunitiesUC)
	logger.Debug("Controllers initialized")

	// データ初期化
	seedData(communityRepo, agentRepo)
	logger.Info("Seed data inserted")

	// ルーティング
	r := router.NewRouter(
		commCtrl,
		diploCtrl,
		simCtrl,
		imageCtrl,
		interferenceCtrl, // 新たに渡す
	)
	logger.Info("Router initialized")

	logger.Info("Starting zousui server on :8080")
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
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
