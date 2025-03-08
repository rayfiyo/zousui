package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/usecase"
	"go.uber.org/zap"
)

type CommunityController struct {
	communityUC *usecase.CommunityUsecase
}

func NewCommunityController(
	uc *usecase.CommunityUsecase,
) *CommunityController {
	zap.L().Debug("Initializing CommunityController")
	return &CommunityController{communityUC: uc}
}

// POST /communities
func (cc *CommunityController) CreateCommunity(
	c *gin.Context,
) {
	logger := zap.L()

	var req struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Population  int    `json:"population"`
		Culture     string `json:"culture"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	newComm := &entity.Community{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Population:  req.Population,
		Culture:     req.Culture,
		UpdatedAt:   time.Now(),
	}
	logger.Debug("Creating community", zap.String("communityID", newComm.ID))

	if err := cc.communityUC.CreateCommunity(c, newComm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Community created", zap.String("communityID", newComm.ID))
	c.JSON(http.StatusOK, gin.H{"message": "created community", "community": newComm})
}

// GET /communities/:id
func (cc *CommunityController) GetCommunity(
	c *gin.Context,
) {
	logger := zap.L()

	id := c.Param("id")
	logger.Debug("Fetching community", zap.String("communityID", id))

	comm, err := cc.communityUC.GetCommunityByID(c, id)
	if err != nil {
		logger.Warn("Community not found", zap.String("communityID", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
		return
	}

	logger.Info("Community fetched", zap.String("communityID", id))
	c.JSON(http.StatusOK, comm)
}

// DELETE /communities/:id
func (cc *CommunityController) DeleteCommunity(
	c *gin.Context,
) {
	logger := zap.L()

	id := c.Param("id")
	logger.Debug("Deleting community", zap.String("communityID", id))

	err := cc.communityUC.DeleteCommunity(c, id)
	if err != nil {
		logger.Error("Failed to delete community", zap.String("communityID", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Community deleted", zap.String("communityID", id))
	c.JSON(http.StatusOK, gin.H{"message": "deleted community", "id": id})
}

// コミュニティ一覧をJSONで返す
func (cc *CommunityController) GetCommunities(
	c *gin.Context,
) {
	logger := zap.L()

	logger.Debug("Fetching all communities")

	communities, err := cc.communityUC.GetAllCommunities(c)
	if err != nil {
		logger.Error("Failed to fetch communities", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Communities fetched", zap.Int("count", len(communities)))
	c.JSON(http.StatusOK, communities)
}
