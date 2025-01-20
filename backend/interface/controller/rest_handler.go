package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
)

// SimulateController: シミュレーションAPIをハンドリング
type SimulateController struct {
	simulateUC *usecase.SimulateCultureEvolutionUsecase
}

// NewSimulateController: コンストラクタ
func NewSimulateController(simUC *usecase.SimulateCultureEvolutionUsecase) *SimulateController {
	return &SimulateController{simulateUC: simUC}
}

// POST /simulate/:communityID
func (sc *SimulateController) Simulate(c *gin.Context) {
	communityID := c.Param("communityID")

	err := sc.simulateUC.Execute(context.Background(), communityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Simulation executed successfully."})
}

// SetupRoutes: ルーティング設定
func (sc *SimulateController) SetupRoutes(r *gin.Engine) {
	r.POST("/simulate/:communityID", sc.Simulate)
}
