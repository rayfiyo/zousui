package entity

import (
	"net/http"
)

// Agent: LLMやAIキャラクター（コミュニティを担当するエージェント）
type Agent struct {
	ID          string
	Name        string
	CommunityID string
	Personality string // 簡易表現：性格や設定など
}

// CultureUpdateRequest: エージェントがコミュニティの文化を更新するためのリクエスト
type CultureUpdateResponse struct {
	NewCulture       string `json:"newCulture"`
	PopulationChange int    `json:"populationChange"`
}

// GeminiLLMGateway: GeminiLLM APIを利用するためのゲートウェイ
type GeminiLLMGateway struct {
	apiKey     string
	httpClient *http.Client
}
