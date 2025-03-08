package gateway

import (
	"context"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type MockLLMGatewayJSON struct{}

// LLMに問い合わせて、文化の変化をJSONで取得する
func (m *MockLLMGatewayJSON) GenerateCultureUpdate(
	ctx context.Context,
	prompt string,
) (string, error) {
	logger := zap.L()

	logger.Debug("Mock GenerateCultureUpdate called", zap.String("prompt", prompt))
	jsonResult := `{
        "newCulture": "踊りを中心にした新たな祭典文化",
        "populationChange": 15
    }`

	logger.Debug("Mock response", zap.String("response", jsonResult))
	fmt.Println("[DEBUG] Prompt to LLM:", prompt)
	fmt.Println("[DEBUG] Return JSON:", jsonResult)
	return jsonResult, nil
}

var _ repository.LLMGateway = (*MockLLMGatewayJSON)(nil)
