package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
	"go.uber.org/zap"
)

// コミュニティ同士の干渉を行うコントローラ
type InterferenceController struct {
	interferenceUC *usecase.SimulateInterferenceBetweenCommunitiesUsecase
}

func NewInterferenceController(
	uc *usecase.SimulateInterferenceBetweenCommunitiesUsecase,
) *InterferenceController {
	zap.L().Debug("Initializing InterferenceController")
	return &InterferenceController{interferenceUC: uc}
}

// POST /simulate/interference?commA=xxx&commB=yyy
func (ic *InterferenceController) SimulateInterferenceBetweenCommunities(
	c *gin.Context,
) {
	logger := zap.L()

	commA := c.Query("commA")
	commB := c.Query("commB")
	logger.Debug("Simulating interference", zap.String("commA", commA), zap.String("commB", commB))

	if commA == "" || commB == "" {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "missing commA or commB parameter"})
		return
	}
	if commA == commB {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "commA and commB must be different"})
		return
	}

	if err := ic.interferenceUC.Execute(c, commA, commB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Interference simulation succeeded",
		zap.String("commA", commA), zap.String("commB", commB))
	c.JSON(http.StatusOK, gin.H{
		"message": "Interference simulation done",
		"commA":   commA,
		"commB":   commB,
	})
}
