package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/usecase"
)

// CommunityController: コミュニティ関連APIをまとめる
type CommunityController struct {
	communityUC *usecase.CommunityUsecase
}

func NewCommunityController(uc *usecase.CommunityUsecase) *CommunityController {
	return &CommunityController{communityUC: uc}
}

// CreateCommunity: POST /communities
func (cc *CommunityController) CreateCommunity(c *gin.Context) {
	var req struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Population  int    `json:"population"`
		Culture     string `json:"culture"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
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
	if err := cc.communityUC.CreateCommunity(c, newComm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "created community", "community": newComm})
}

// GetCommunity: GET /communities/:id
func (cc *CommunityController) GetCommunity(c *gin.Context) {
	id := c.Param("id")
	comm, err := cc.communityUC.GetCommunityByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
		return
	}
	c.JSON(http.StatusOK, comm)
}

// DeleteCommunity: DELETE /communities/:id
func (cc *CommunityController) DeleteCommunity(c *gin.Context) {
	id := c.Param("id")
	err := cc.communityUC.DeleteCommunity(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted community", "id": id})
}

// GetCommunities: コミュニティ一覧をJSONで返す
func (cc *CommunityController) GetCommunities(c *gin.Context) {
	communities, err := cc.communityUC.GetAllCommunities(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// JSON出力
	c.JSON(http.StatusOK, communities)
}
