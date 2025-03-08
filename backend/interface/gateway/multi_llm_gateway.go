package gateway

import (
	"context"
	"math/rand"

	"github.com/rayfiyo/zousui/backend/domain/repository"
	"go.uber.org/zap"
)

type MultiLLMGateway struct {
	subGateways []repository.LLMGateway
}

func NewMultiLLMGateway(
	subGateways ...repository.LLMGateway,
) *MultiLLMGateway {
	return &MultiLLMGateway{
		subGateways: subGateways,
	}
}

// 複数のLLMへの問い合わせ結果から「ランダムで1つ」選んで返す例
func (m *MultiLLMGateway) GenerateCultureUpdate(
	ctx context.Context,
	prompt string,
) (string, error) {
	logger := zap.L()

	if len(m.subGateways) == 0 {
		logger.Error("No sub gateways available in MultiLLMGateway")
		return "", nil
	}

	r := rand.New(rand.NewSource(42))
	gwIndex := r.Intn(len(m.subGateways))
	chosen := m.subGateways[gwIndex]

	logger.Debug("Selected sub gateway", zap.Int("index", gwIndex))
	return chosen.GenerateCultureUpdate(ctx, prompt)
}

var _ repository.LLMGateway = (*MultiLLMGateway)(nil)
