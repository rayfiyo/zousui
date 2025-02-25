package gateway

import (
	"context"
	"math/rand"
	"time"

	"github.com/rayfiyo/zousui/backend/domain/repository"
)

// MultiLLMGateway: 複数のLLMゲートウェイをまとめて扱うゲートウェイ
type MultiLLMGateway struct {
	subGateways []repository.LLMGateway
}

// NewMultiLLMGateway: コンストラクタ
// サブゲートウェイを可変長で渡す
func NewMultiLLMGateway(subGateways ...repository.LLMGateway) *MultiLLMGateway {
	return &MultiLLMGateway{
		subGateways: subGateways,
	}
}

// GenerateCultureUpdate: 複数のLLMへの問い合わせ結果から「ランダムで1つ」選んで返す例
func (m *MultiLLMGateway) GenerateCultureUpdate(ctx context.Context, prompt string) (string, error) {
	if len(m.subGateways) == 0 {
		// サブゲートウェイが無い場合のエラー
		return "", nil
	}

	// ランダムなサブゲートウェイを選ぶ
	rand.Seed(time.Now().UnixNano())
	gwIndex := rand.Intn(len(m.subGateways))
	chosen := m.subGateways[gwIndex]

	// 選んだゲートウェイで生成
	return chosen.GenerateCultureUpdate(ctx, prompt)
}

// Compile-time interface check
var _ repository.LLMGateway = (*MultiLLMGateway)(nil)
