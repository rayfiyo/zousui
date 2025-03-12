package entity

import "time"

// シミュレーションの結果を表します。
type SimulationResult struct {
	ID          string    // 一意のID（例：UUID）
	Type        string    // "diplomacy", "interference" など
	Communities []string  // 関連するコミュニティIDの一覧
	ResultJSON  string    // シミュレーション結果のJSON文字列
	CreatedAt   time.Time // 実行日時
}
