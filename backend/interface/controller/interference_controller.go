package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
	"go.uber.org/zap"
)

// リクエストボディの構造体
type InterferenceRequest struct {
	CommA     string `json:"commA"`
	CommB     string `json:"commB"`
	UserInput string `json:"userInput"`
}

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

	var req InterferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.CommA == "" || req.CommB == "" || req.CommA == req.CommB {
		logger.Error("commA and commB must be provided and different")
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "commA and commB must be provided and different"})
		return
	}

	// Usecase実行
	if err := ic.interferenceUC.Execute(
		c, req.CommA, req.CommB, req.UserInput,
	); err != nil {
		logger.Error("failed to execute interference simulation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功応答
	c.JSON(http.StatusOK, gin.H{
		"message": "Interference simulation done",
		"commA":   req.CommA,
		"commB":   req.CommB,
	})
}
