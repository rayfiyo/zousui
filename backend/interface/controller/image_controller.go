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
	"go.uber.org/zap"
)

// 文化(Culture)テキストを画像生成APIに渡して画像を返却するコントローラ
type ImageController struct {
	communityUC usecase.CommunityUsecase
}

func NewImageController(
	cu usecase.CommunityUsecase,
) *ImageController {
	zap.L().Debug("Initializing ImageController")
	return &ImageController{
		communityUC: cu,
	}
}

// リクエストボディ用の構造体
type GenerateImageRequest struct {
	Style string `json:"style,omitempty"`
}

// OpenAIの画像生成APIレスポンス
type imageResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// POST /communities/:communityID/generateImage
func (ic *ImageController) GenerateImage(
	c *gin.Context,
) {
	logger := zap.L()

	communityID := c.Param("communityID")
	logger.Debug("GenerateImage called", zap.String("communityID", communityID))

	comm, err := ic.communityUC.GetCommunityByID(c, communityID)
	if err != nil || comm == nil {
		logger.Warn("Community not found for image generation",
			zap.String("communityID", communityID))
		c.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
		return
	}

	var req GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("No request body for image generation or invalid JSON",
			zap.Error(err))
	}

	if config.OpenAIAPIKEY == "" {
		logger.Error("OPENAI_API_KEY is not set")
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "OPENAI_API_KEY environment variable is not set"})
		return
	}

	prompt := fmt.Sprintf("A digital art representing the culture: %q", comm.Culture)
	if req.Style != "" {
		prompt += fmt.Sprintf(", style: %s", req.Style)
	}
	logger.Debug("Image generation prompt", zap.String("prompt", prompt))

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  consts.DALLEModel,
		"prompt": prompt,
		"n":      1,
		"size":   consts.ImageSize,
	})
	if err != nil {
		logger.Error("Failed to marshal image request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpReq, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, consts.DALLEEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		logger.Error("Failed to create HTTP request for image generation",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+config.OpenAIAPIKEY)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Error("Failed to execute image generation HTTP request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody bytes.Buffer
		errBody.ReadFrom(resp.Body)
		logger.Error("OpenAI returned non-OK status", zap.Int("status",
			resp.StatusCode), zap.String("response", errBody.String()))
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("OpenAI error: %s", errBody.String())})
		return
	}

	var imageResp imageResponse
	if err := json.NewDecoder(resp.Body).Decode(&imageResp); err != nil {
		logger.Error("Failed to decode image response", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(imageResp.Data) == 0 {
		logger.Error("No image data returned")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no image data returned"})
		return
	}

	imageURL := imageResp.Data[0].URL
	logger.Debug("Downloading image", zap.String("imageURL", imageURL))
	imgResp, err := http.Get(imageURL)
	if err != nil {
		logger.Error("Failed to download image", zap.Error(err))
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
		logger.Error("Unsupported image format", zap.String("contentType", contentType))
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Unsupported image format: " + contentType})
		return
	}
	if err != nil {
		logger.Error("Failed to decode image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resized := resize.Resize(300, 0, decoded, resize.Lanczos3)
	logger.Info("Image processed", zap.String("communityID", communityID))

	c.Writer.Header().Set("Content-Type", contentType)
	switch contentType {
	case "image/jpeg":
		jpeg.Encode(c.Writer, resized, nil)
	case "image/png":
		png.Encode(c.Writer, resized)
	}
}
