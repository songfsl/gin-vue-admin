package dto

// ReviewListResponse レビュー一覧APIのルートレスポンス
type ReviewListResponse struct {
	Summary    *ReviewSummary `json:"summary"`    // レビュー集計情報 (Nullable: 商品にレビューがない場合)
	Reviews    []ReviewInfo   `json:"reviews"`    // レビューリスト
	Pagination PaginationInfo `json:"pagination"` // ページネーション情報
}

// ReviewSummary レビュー集計情報
type ReviewSummary struct {
	AverageRating float64 `json:"average_rating" gorm:"column:average_rating"`
	ReviewCount   int     `json:"review_count" gorm:"column:review_count"`
	Rating1Count  int     `json:"rating_1_count" gorm:"column:rating_1_count"`
	Rating2Count  int     `json:"rating_2_count" gorm:"column:rating_2_count"`
	Rating3Count  int     `json:"rating_3_count" gorm:"column:rating_3_count"`
	Rating4Count  int     `json:"rating_4_count" gorm:"column:rating_4_count"`
	Rating5Count  int     `json:"rating_5_count" gorm:"column:rating_5_count"`
}

// ReviewInfo 個々のレビュー情報 (修正: image_urls, helpful_count 追加)
type ReviewInfo struct {
	ID                    int64    `json:"id"`                   // レビューID
	Nickname              string   `json:"nickname"`             // ニックネーム
	Rating                int      `json:"rating"`               // 評価 (1-5)
	Title                 *string  `json:"title,omitempty"`      // タイトル (Nullable)
	Comment               string   `json:"comment"`              // 本文
	CreatedAtFormatted    string   `json:"-"`                    // 表示用投稿日時 (例: "2023年10月26日")
	CreatedAtFormattedStr string   `json:"created_at_formatted"` // 表示用投稿日時 (例: "2023年10月26日")
	ImageUrls             string   `json:"-"`                    // 用于接收GROUP_CONCAT的字符串，不对外输出
	RealImageUrls         []string `json:"image_urls,omitempty"` // 对外输出的字段
	HelpfulCount          int      `json:"helpful_count"`        // ★参考になった数
	// IsHelpfulByUser   *bool    `json:"is_helpful_by_user,omitempty"` // ★(オプション) ログインユーザーが参考になったを押したか (Nullable)
}

// PaginationInfo ページネーション情報 (変更なし)
type ReviewPaginationInfo struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalCount  int `json:"total_count"`
	TotalPages  int `json:"total_pages"`
}

// ErrorResponse エラーレスポンス構造 (共通)
type ReviewErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
type ReviewErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}
