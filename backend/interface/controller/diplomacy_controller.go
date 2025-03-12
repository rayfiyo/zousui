package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
	"go.uber.org/zap"
)

type DiplomacyController struct {
	diploUC *usecase.DiplomacyUsecase
}

func NewDiplomacyController(
	uc *usecase.DiplomacyUsecase,
) *DiplomacyController {
	zap.L().Debug("Initializing DiplomacyController")
	return &DiplomacyController{diploUC: uc}
}

// POST /simulate/diplomacy?commA=...&commB=...
func (dc *DiplomacyController) SimulateDiplomacy(
	c *gin.Context,
) {
	logger := zap.L()

	commA := c.Query("commA")
	commB := c.Query("commB")
	logger.Debug("Simulating diplomacy", zap.String("commA", commA),
		zap.String("commB", commB))

	if commA == "" || commB == "" {
		logger.Warn("Missing parameters for diplomacy simulation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "commA and commB are required"})
		return
	}

	if err := dc.diploUC.ExecuteDiplomacy(c, commA, commB); err != nil {
		logger.Error("Diplomacy simulation failed", zap.Error(err),
			zap.String("commA", commA), zap.String("commB", commB))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Diplomacy simulation succeeded", zap.String("commA", commA),
		zap.String("commB", commB))
	c.JSON(http.StatusOK, gin.H{"message": "Diplomacy simulation done"})
}
