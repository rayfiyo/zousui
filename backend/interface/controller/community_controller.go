package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
)

// CommunityController: コミュニティ関連APIをまとめる
type CommunityController struct {
	communityUC *usecase.CommunityUsecase
}

func NewCommunityController(uc *usecase.CommunityUsecase) *CommunityController {
	return &CommunityController{communityUC: uc}
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
