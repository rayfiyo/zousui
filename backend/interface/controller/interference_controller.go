package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
)

// コミュニティ同士の干渉を行うコントローラ
type InterferenceController struct {
	interferenceUC *usecase.SimulateInterferenceBetweenCommunitiesUsecase
}

// コンストラクタ
func NewInterferenceController(uc *usecase.SimulateInterferenceBetweenCommunitiesUsecase) *InterferenceController {
	return &InterferenceController{interferenceUC: uc}
}

// POST /simulate/interference?commA=xxx&commB=yyy
func (ic *InterferenceController) SimulateInterferenceBetweenCommunities(c *gin.Context) {
	commA := c.Query("commA")
	commB := c.Query("commB")
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

	// Usecase実行
	err := ic.interferenceUC.Execute(c, commA, commB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功応答
	c.JSON(http.StatusOK, gin.H{
		"message": "Interference simulation done",
		"commA":   commA,
		"commB":   commB,
	})
}
