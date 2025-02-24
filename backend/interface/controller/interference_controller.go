package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
)

// InterferenceController: 複数LLM干渉シミュレーション用のコントローラ
type InterferenceController struct {
	interferenceUC *usecase.SimulateInterferenceUsecase
}

// NewInterferenceController: コンストラクタ
func NewInterferenceController(iuc *usecase.SimulateInterferenceUsecase) *InterferenceController {
	return &InterferenceController{interferenceUC: iuc}
}

// SimulateInterference: POST /simulate/interference/:communityID
func (ic *InterferenceController) SimulateInterference(c *gin.Context) {
	communityID := c.Param("communityID")
	if communityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid communityID"})
		return
	}

	err := ic.interferenceUC.Execute(context.Background(), communityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interference simulation succeeded."})
}
