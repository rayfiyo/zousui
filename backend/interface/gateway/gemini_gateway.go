package gateway

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	// "github.com/rayfiyo/zousui/backend/usecase"
	"github.com/rayfiyo/zousui/backend/utils"
	"github.com/rayfiyo/zousui/backend/utils/config"
	"google.golang.org/api/option"
)

type GeminiLLMGateway struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiLLMGateway: コンストラクタ
func NewGeminiLLMGateway(ctx context.Context) (*GeminiLLMGateway, error) {
	if config.GeminiAPIKEY == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GeminiAPIKEY))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	// [TODO]
	model := client.GenerativeModel("gemini")

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(utils.SpecifyingResponseFormat),
		},
	}

	return &GeminiLLMGateway{
		client: client,
		model:  model,
	}, nil
}

// GenerateCultureUpdate: LLMGatewayインタフェース
// ここではストリーミングなしで結果をまとめて取得する例
func (g *GeminiLLMGateway) GenerateCultureUpdate(ctx context.Context, prompt string) (string, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	//  [TODO] JSON にパースして返す

	return fmt.Sprintln(resp.Candidates[0].Content.Parts), nil
}

// var _ usecase.LLMGateway = (*GeminiLLMGateway)(nil)
