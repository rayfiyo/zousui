package gateway

import (
	"context"

	"github.com/rayfiyo/zousui/backend/usecase"
)

// MockLLMGateway: PoC用に、LLMへの問い合わせをモック化(実際にはOpenAIや他APIへHTTPリクエスト)
type MockLLMGateway struct{}

// GenerateCultureUpdate: ダミー実装
func (m *MockLLMGateway) GenerateCultureUpdate(ctx context.Context, prompt string) (string, error) {
	// 本来は外部APIへHTTPリクエストなど行う
	return "【新しい文化】神秘の踊りを中心とした祭典を毎週開催し、踊り子がリーダーとして社会を動かす制度", nil
}

// インタフェースが正しく実装されているかコンパイル時チェック
var _ usecase.LLMGateway = (*MockLLMGateway)(nil)
