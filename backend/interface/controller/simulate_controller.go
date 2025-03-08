package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
	"go.uber.org/zap"
)

type SimulateController struct {
	simulateUC *usecase.SimulateCultureEvolutionUsecase
}

func NewSimulateController(
	simUC *usecase.SimulateCultureEvolutionUsecase,
) *SimulateController {
	zap.L().Debug("Initializing SimulateController")
	return &SimulateController{simulateUC: simUC}
}

// POST /simulate/:communityID
func (sc *SimulateController) Simulate(
	c *gin.Context,
) {
	logger := zap.L()

	communityID := c.Param("communityID")
	logger.Debug("Simulate called", zap.String("communityID", communityID))

	if err := sc.simulateUC.Execute(context.Background(), communityID); err != nil {
		logger.Error("Simulation failed",
			zap.String("communityID", communityID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Simulation succeeded", zap.String("communityID", communityID))
	c.JSON(http.StatusOK, gin.H{"message": "Simulation executed successfully."})
}

// ルーティング設定
func (sc *SimulateController) SetupRoutes(r *gin.Engine) {
	r.POST("/simulate/:communityID", sc.Simulate)
}
