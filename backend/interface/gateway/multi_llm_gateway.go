package gateway

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/rayfiyo/zousui/backend/domain/repository"
	"github.com/rayfiyo/zousui/backend/utils/consts"
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

// テキストからランダムな単語を抽出する
func extractRandomWord(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return words[r.Intn(len(words))]
}

// mergeResponses は複数の回答を改行区切りでマージする
func mergeResponses(responses []string) string {
	return strings.Join(responses, "\n")
}

// 各サブゲートウェイへ並列問い合わせし、各回答をマージした上で、
// ユーザー入力を加えた集約プロンプトを作成し、再問い合わせする
func (m *MultiLLMGateway) GenerateCultureUpdate(
	ctx context.Context,
	prompt string,
	userInput string,
) (string, error) {
	logger := zap.L()

	if len(m.subGateways) == 0 {
		logger.Error("No sub gateways available in MultiLLMGateway")
		return "", errors.New("no sub gateways available")
	}

	// 並列呼び出し
	type resp struct {
		output string
		err    error
	}
	responses := make([]resp, len(m.subGateways))
	var wg sync.WaitGroup
	wg.Add(len(m.subGateways))
	for i, gw := range m.subGateways {
		go func(i int, gw repository.LLMGateway) {
			defer wg.Done()
			out, err := gw.GenerateCultureUpdate(ctx, prompt, userInput)
			responses[i] = resp{output: out, err: err}
		}(i, gw)
	}
	wg.Wait()

	valid := []string{}
	for _, r := range responses {
		if r.err == nil && r.output != "" {
			valid = append(valid, r.output)
		}
	}
	if len(valid) == 0 {
		logger.Error("All sub gateway calls failed")
		return "", errors.New("all sub gateway calls failed")
	}

	// 複数の回答があればマージ
	var merged string
	if len(valid) == 1 {
		merged = valid[0]
	} else {
		merged = mergeResponses(valid)
	}

	// マージ結果からランダムな単語を抽出（追加のインスピレーションとして利用）
	randomWord := extractRandomWord(valid[rand.Intn(len(valid))])
	if randomWord == "" {
		randomWord = "革新"
	}

	// 集約プロンプトの生成
	aggregatedPrompt := consts.AggregatedPromptHeader + "\n" + "キーワード: " +
		randomWord + "追加情報: " + userInput + "複数のアイデア: " + merged
	logger.Debug("Aggregated prompt", zap.String("aggregatedPrompt", aggregatedPrompt))

	// 先頭のサブゲートウェイに再問い合わせして最終結果を得る
	finalResponse, err := m.subGateways[0].GenerateCultureUpdate(
		ctx, aggregatedPrompt, userInput)
	if err != nil {
		logger.Error("Final aggregated call failed", zap.Error(err))
		return "", err
	}
	logger.Debug("Final aggregated response", zap.String("response", finalResponse))
	return finalResponse, nil
}

var _ repository.LLMGateway = (*MultiLLMGateway)(nil)
