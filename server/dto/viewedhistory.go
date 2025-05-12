package dto

// ★ DTO構造体名をバージョンアップ等で変更推奨 (例: V1_1) ★

// ViewedSKUListResponse 閲覧履歴リストAPIのルートレスポンス (修正版)
type ViewedSKUListResponse struct {
	History    []ViewedSKUInfo `json:"history"`    // 閲覧履歴SKUリスト
	Pagination PaginationInfo  `json:"pagination"` // ページネーション情報
}

// ViewedSKUInfoV1_1 個々の閲覧履歴SKU情報 (修正版)
type ViewedSKUInfo struct {
	SkuID               string             `json:"sku_id"`                   // SKU ID
	ProductID           string             `json:"product_id"`               // 商品ID
	ProductName         string             `json:"product_name"`             // 商品名 (省略後)
	ProductCode         string             `json:"product_code,omitempty"`   // 商品コード
	PriceRangeFormatted string             `json:"price_range_formatted"`    // ★ 商品の価格帯文字列
	PrimaryImage        *ImageInfo         `json:"primary_image,omitempty"`  // サムネイル画像推奨 (Nullable)
	ReviewSummary       *ReviewSummaryInfo `json:"review_summary,omitempty"` // ★ 商品のレビュー集計情報 (Nullable)
	ViewedAtFormatted   string             `json:"viewed_at_formatted"`      // 最終閲覧日時 (表示用)
	// Price             *PriceInfo      `json:"price"`                     // ★削除
	// Attributes        []AttributeInfo `json:"attributes"`                // ★削除
}

// ReviewSummaryInfo レビュー集計情報のDTO (関連商品APIと共通化可能)
/*
type ReviewSummaryInfo struct {
	AverageRating float64 `json:"average_rating"` // 平均評価
	ReviewCount   int     `json:"review_count"`   // レビュー件数
}

// --- 以下のDTOは他のAPIと共通化 ---

// ImageInfo 画像情報
type ImageInfo struct {
	ID      int     `json:"id"`
	URL     string  `json:"url"` // サムネイルURL想定
	AltText *string `json:"alt_text,omitempty"`
}

// PaginationInfo ページネーション情報
type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalCount  int `json:"total_count"` // 閲覧履歴の総件数
	TotalPages  int `json:"total_pages"`
}

// ErrorResponse エラーレスポンス構造 (共通)
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}
*/
