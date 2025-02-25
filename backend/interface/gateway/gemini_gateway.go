package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/rayfiyo/zousui/backend/domain/repository"
	"github.com/rayfiyo/zousui/backend/utils/config"
	"github.com/rayfiyo/zousui/backend/utils/consts"
	"google.golang.org/api/option"
)

type GeminiLLMGateway struct {
	Client *genai.Client
	Model  *genai.GenerativeModel
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
	if client == nil {
		return nil, fmt.Errorf("gemini client is nil")
	}

	model := client.GenerativeModel(consts.GeminiModel)

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(consts.SpecifyingResponseFormat),
		},
	}

	return &GeminiLLMGateway{
		Client: client,
		Model:  model,
	}, nil
}

// GenerateCultureUpdate: LLMGatewayインタフェース
// ここではストリーミングなしで結果をまとめて取得する例
func (g *GeminiLLMGateway) GenerateCultureUpdate(ctx context.Context, prompt string) (string, error) {
	respRaw, err := g.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// 装飾の削除
	resp := strings.ReplaceAll(
		fmt.Sprintln(respRaw.Candidates[0].Content.Parts), "```", "")
	resp = strings.ReplaceAll(resp, "json", "")
	resp = strings.ReplaceAll(resp, "[", "")
	resp = strings.ReplaceAll(resp, "]", "")

	return resp, nil
}

// インタフェースが正しく実装されているかコンパイル時チェック
var _ repository.LLMGateway = (*GeminiLLMGateway)(nil)
