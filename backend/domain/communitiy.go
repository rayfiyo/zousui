package domain

import "time"

// Community: 仮想社会（コミュニティ）を表すドメインモデル
type Community struct {
	ID          string
	Name        string
	Description string
	Population  int
	Culture     string
	UpdatedAt   time.Time
}

// UpdateCulture: コミュニティの文化情報を更新するドメインロジック
func (c *Community) UpdateCulture(newCulture string) {
	c.Culture = newCulture
	c.UpdatedAt = time.Now()
}
