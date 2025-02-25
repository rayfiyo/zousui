package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/rayfiyo/zousui/backend/usecase"
	"github.com/rayfiyo/zousui/backend/utils/config"
	"github.com/rayfiyo/zousui/backend/utils/consts"
)

// ImageController: 文化(Culture)テキストを画像生成APIに渡して画像を返却するコントローラ
type ImageController struct {
	communityUC usecase.CommunityUsecase
}

func NewImageController(cu usecase.CommunityUsecase) *ImageController {
	return &ImageController{
		communityUC: cu,
	}
}

// GenerateImageRequest: リクエストボディ用の構造体（将来的な拡張を想定）
type GenerateImageRequest struct {
	Style string `json:"style,omitempty"`
}

// imageResponse: OpenAIの画像生成APIレスポンス(簡易)
type imageResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// POST /communities/:communityID/generateImage
func (ic *ImageController) GenerateImage(c *gin.Context) {
	// 1. communityID を取得
	communityID := c.Param("communityID")

	// 2. DBからコミュニティ情報を取得 (ユースケースを呼ぶ)
	comm, err := ic.communityUC.GetCommunityByID(c, communityID)
	if err != nil || comm == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
		return
	}

	// 3. bodyをパース(拡張用)。無ければ使わずにOK
	var req GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// body無しの場合は無視 or エラー
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// return
	}

	// 4. 画像生成APIキーを取得
	if config.OpenAIAPIKEY == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OPENAI_API_KEY environment variable is not set"})
		return
	}

	// 5. 生成用のプロンプトを組み立て
	//    ここではコミュニティの Culture を使う例
	//    例: "A beautiful digital art representing the culture: '～～'"
	prompt := fmt.Sprintf("A digital art representing the culture: %q", comm.Culture)
	if req.Style != "" {
		// styleがあれば追加
		prompt += fmt.Sprintf(", style: %s", req.Style)
	}

	// 6. OpenAIの画像生成APIへリクエスト
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  consts.DALLEModel,
		"prompt": prompt,
		"n":      1,
		"size":   consts.ImageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpReq, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, consts.DALLEEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+config.OpenAIAPIKEY)

	// 7. リクエスト送信
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// エラーボディを読み込んで返す
		var errBody bytes.Buffer
		errBody.ReadFrom(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("OpenAI error: %s", errBody.String())})
		return
	}

	var imageResp imageResponse
	if err := json.NewDecoder(resp.Body).Decode(&imageResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(imageResp.Data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no image data returned"})
		return
	}

	// 8. 実際の画像URLからダウンロード
	imageURL := imageResp.Data[0].URL
	imgResp, err := http.Get(imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer imgResp.Body.Close()

	contentType := imgResp.Header.Get("Content-Type")
	var decoded image.Image
	switch contentType {
	case "image/jpeg":
		decoded, err = jpeg.Decode(imgResp.Body)
	case "image/png":
		decoded, err = png.Decode(imgResp.Body)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsupported image format: " + contentType})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 9. リサイズ (幅300px, 高さアスペクト比維持)
	resized := resize.Resize(300, 0, decoded, resize.Lanczos3)

	// 10. 画像としてレスポンス (mimeそのまま返す)
	c.Writer.Header().Set("Content-Type", contentType)
	switch contentType {
	case "image/jpeg":
		jpeg.Encode(c.Writer, resized, nil)
	case "image/png":
		png.Encode(c.Writer, resized)
	}
}
