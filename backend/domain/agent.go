package domain

// Agent: LLMやAIキャラクター（コミュニティを担当するエージェント）の例
type Agent struct {
	ID          string
	Name        string
	CommunityID string
	Personality string // 簡易表現：性格や設定など
}
