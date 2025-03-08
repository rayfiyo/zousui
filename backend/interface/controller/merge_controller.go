package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/zousui/backend/usecase"
)

// コミュニティ合体エンドポイントを定義
type MergeController struct {
	mergeUC *usecase.MergeCommunitiesUsecase
}

// コンストラクタ
func NewMergeController(mu *usecase.MergeCommunitiesUsecase) *MergeController {
	return &MergeController{mergeUC: mu}
}

// 合体時のリクエストボディ
type requestBody struct {
	CommA       string `json:"commA"`
	CommB       string `json:"commB"`
	NewCommID   string `json:"newID"`
	NewCommName string `json:"newName"`
}

// POST /communities/merge
func (mc *MergeController) MergeCommunities(c *gin.Context) {
	var body requestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	// リクエストバリデーション
	if body.CommA == "" || body.CommB == "" || body.NewCommID == "" || body.NewCommName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "commA, commB, newID, newName are required"})
		return
	}

	// ユースケース呼び出し
	merged, err := mc.mergeUC.Merge(
		context.Background(),
		body.CommA,
		body.CommB,
		body.NewCommID,
		body.NewCommName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功時に新コミュニティ情報を返す
	c.JSON(http.StatusOK, gin.H{
		"message":   "Communities merged successfully",
		"community": merged,
	})
}
