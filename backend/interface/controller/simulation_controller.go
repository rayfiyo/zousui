package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/domain/repository"
)

type SimulationController struct {
	simulationRepo repository.SimulationRepository
}

func NewSimulationController(repo repository.SimulationRepository) *SimulationController {
	return &SimulationController{simulationRepo: repo}
}

func (sc *SimulationController) GetSimulationHistory(c *gin.Context) {
	simulations, err := sc.simulationRepo.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, simulations)
}
