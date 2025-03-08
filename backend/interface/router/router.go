package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/interface/controller"
	"go.uber.org/zap"
)

func NewRouter(
	commCtrl *controller.CommunityController,
	diploCtrl *controller.DiplomacyController,
	simCtrl *controller.SimulateController,
	imageCtrl *controller.ImageController,
	interferenceCtrl *controller.InterferenceController,
) *gin.Engine {
	logger := zap.L()

	logger.Debug("Initializing router")

	// Gin
	r := gin.Default()
	r.Use(cors.Default())

	// コミュニティ一覧 + CRUD
	r.GET("/communities", commCtrl.GetCommunities)
	r.GET("/communities/:id", commCtrl.GetCommunity)
	r.POST("/communities", commCtrl.CreateCommunity)
	r.DELETE("/communities/:id", commCtrl.DeleteCommunity)

	// 外交シミュレーション
	r.POST("/simulate/diplomacy", diploCtrl.SimulateDiplomacy)

	// シミュレーション実行
	r.POST("/simulate/:communityID", simCtrl.Simulate)

	// 画像生成API
	r.POST("/communities/:communityID/generateImage", imageCtrl.GenerateImage)

	// 干渉シミュレーション (コミュニティAとB)
	r.POST("/simulate/interference",
		interferenceCtrl.SimulateInterferenceBetweenCommunities)

	logger.Info("Router initialized")
	return r
}
