package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rayfiyo/zousui/backend/domain/entity"
	"github.com/rayfiyo/zousui/backend/domain/repository"
)

// 2つのコミュニティを合体(新コミュニティとして生成)するユースケース
type MergeCommunitiesUsecase struct {
	communityRepo repository.CommunityRepository
	llmGateway    repository.LLMGateway
}

// コンストラクタ
func NewMergeCommunitiesUsecase(
	cr repository.CommunityRepository,
	lg repository.LLMGateway,
) *MergeCommunitiesUsecase {
	return &MergeCommunitiesUsecase{
		communityRepo: cr,
		llmGateway:    lg,
	}
}

// LLMからの文化合体用レスポンス(JSON)
// 例: { "mergedCulture": "新しい融合文化" }
type mergeCultureResponse struct {
	MergedCulture string `json:"mergedCulture"`
}

// commA と commB を合体し、1つの新コミュニティを生成する
//
// mergeReq:
//   - commAID: 合体元コミュニティAのID
//   - commBID: 合体元コミュニティBのID
//   - newCommID: 生成される新コミュニティのID
//   - newCommName: 生成される新コミュニティの名前
func (m *MergeCommunitiesUsecase) Merge(ctx context.Context,
	commAID, commBID, newCommID, newCommName string,
) (*entity.Community, error) {
	// 1. コミュニティA, Bを取得
	commA, err := m.communityRepo.GetByID(ctx, commAID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community A: %w", err)
	}
	commB, err := m.communityRepo.GetByID(ctx, commBID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community B: %w", err)
	}

	// 2. LLMプロンプト用のテキストを組み立て
	prompt := fmt.Sprintf(`
        コミュニティAの文化: %q
        コミュニティBの文化: %q
        これら2つの文化を融合して、新たな文化を作ってください。
        出力は必ず以下のJSON形式にしてください:
        {
            "mergedCulture": "～～～"
        }`, commA.Culture, commB.Culture)

	// 3. LLM呼び出し
	llmResp, err := m.llmGateway.GenerateCultureUpdate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM for merging: %w", err)
	}

	// 4. JSONパース
	var mergeResp mergeCultureResponse
	if err := json.Unmarshal([]byte(llmResp), &mergeResp); err != nil {
		// JSONパースに失敗した場合は、単純にテキストをマージ文化として扱う
		mergeResp.MergedCulture = commA.Culture + " + " + commB.Culture
	}

	// 5. 新しいコミュニティを構築
	mergedPopulation := commA.Population + commB.Population
	if mergedPopulation < 0 {
		mergedPopulation = 0
	}

	newComm := &entity.Community{
		ID:          newCommID,
		Name:        newCommName,
		Description: fmt.Sprintf("Merged from %s and %s", commA.Name, commB.Name),
		Population:  mergedPopulation,
		Culture:     mergeResp.MergedCulture,
		UpdatedAt:   time.Now(),
	}

	// 6. 新コミュニティを保存
	if err := m.communityRepo.Save(ctx, newComm); err != nil {
		return nil, fmt.Errorf("failed to save new merged community: %w", err)
	}

	// 7. 旧コミュニティを削除（「合体」なので、残さない想定）
	//    必要に応じて、ここをコメントアウトして「残す」実装に変えてもOK
	if err := m.communityRepo.Delete(ctx, commAID); err != nil {
		return nil, fmt.Errorf("failed to delete community A: %w", err)
	}
	if err := m.communityRepo.Delete(ctx, commBID); err != nil {
		return nil, fmt.Errorf("failed to delete community B: %w", err)
	}

	// 結果を返却
	return newComm, nil
}
