package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"github.com/rayfiyo/zousui/backend/utils/config"
	"github.com/rayfiyo/zousui/backend/utils/consts"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type GeminiLLMGateway struct {
	Client *genai.Client
	Model  *genai.GenerativeModel
}

func NewGeminiLLMGateway(
	ctx context.Context,
) (*GeminiLLMGateway, error) {
	logger := zap.L()

	if config.GeminiAPIKEY == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GeminiAPIKEY))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("gemini client is nil")
	}

	model := client.GenerativeModel(consts.GeminiModel)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(consts.SpecifyingResponseFormat),
		},
	}

	logger.Info("GeminiLLMGateway created")
	return &GeminiLLMGateway{
		Client: client,
		Model:  model,
	}, nil
}

// LLMGatewayインタフェース
func (g *GeminiLLMGateway) GenerateCultureUpdate(
	ctx context.Context,
	prompt string,
	userInput string,
) (string, error) {
	logger := zap.L()

	logger.Debug("Generating culture update with Gemini", zap.String("prompt", prompt))
	respRaw, err := g.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		logger.Error("Failed to generate content", zap.Error(err))
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// 装飾の削除
	resp := strings.ReplaceAll(fmt.Sprintln(respRaw.Candidates[0].Content.Parts),
		"```", "")
	resp = strings.ReplaceAll(resp, "json", "")
	resp = strings.ReplaceAll(resp, "[", "")
	resp = strings.ReplaceAll(resp, "]", "")

	logger.Debug("Generated culture update", zap.String("response", resp))
	return resp, nil
}

var _ repository.LLMGateway = (*GeminiLLMGateway)(nil)
