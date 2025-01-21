package gateway

import (
	"context"
	"fmt"

	"github.com/rayfiyo/zousui/backend/domain/repository"
)

// MockLLMGatewayJSON: LLMへの問い合わせ結果のJSONをモック化
type MockLLMGatewayJSON struct{}

// GenerateCultureUpdate: LLMに問い合わせて、文化の変化をJSONで取得する
func (m *MockLLMGatewayJSON) GenerateCultureUpdate(ctx context.Context, prompt string) (string, error) {
	// ここではダミーでJSONを返す
	// [TODO] OpenAI APIにHTTPリクエストして、そこからのレスポンスを返す
	jsonResult := `{
        "newCulture": "踊りを中心にした新たな祭典文化",
        "populationChange": 15
    }`
	fmt.Println("[DEBUG] Prompt to LLM:", prompt)
	fmt.Println("[DEBUG] Return JSON:", jsonResult)
	return jsonResult, nil
}

// インタフェースが正しく実装されているかコンパイル時チェック
var _ repository.LLMGateway = (*MockLLMGatewayJSON)(nil)
