package gateway

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"github.com/rayfiyo/zousui/backend/utils/config"
	"github.com/rayfiyo/zousui/backend/utils/consts"
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
			genai.Text(consts.SpecifyingResponseFormat),
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

// インタフェースが正しく実装されているかコンパイル時チェック
var _ repository.LLMGateway = (*GeminiLLMGateway)(nil)
