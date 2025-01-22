package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
)

type DiplomacyController struct {
	diploUC *usecase.DiplomacyUsecase
}

func NewDiplomacyController(uc *usecase.DiplomacyUsecase) *DiplomacyController {
	return &DiplomacyController{diploUC: uc}
}

// SimulateDiplomacy: POST /simulate/diplomacy?commA=...&commB=...
func (dc *DiplomacyController) SimulateDiplomacy(c *gin.Context) {
	commA := c.Query("commA")
	commB := c.Query("commB")
	if commA == "" || commB == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "commA and commB are required"})
		return
	}
	if err := dc.diploUC.ExecuteDiplomacy(c, commA, commB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Diplomacy simulation done"})
}
