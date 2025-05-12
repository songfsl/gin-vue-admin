package dto

// PriceInfo等で必要になる可能性

// FavoriteSKUListResponse お気に入りSKUリストAPIのルートレスポンス
type FavoriteSKUListResponse struct {
	Favorites  []FavoriteSKUInfo `json:"favorites"`  // お気に入りSKUリスト
	Pagination PaginationInfo    `json:"pagination"` // ページネーション情報
}

// FavoriteSKUInfo 個々のお気に入りSKU情報
type FavoriteSKUInfo struct {
	SkuID       string     `json:"sku_id"`                 // SKU ID
	ProductID   string     `json:"product_id"`             // 商品ID
	ProductName string     `json:"product_name"`           // 商品名 (省略後)
	ProductCode string     `json:"product_code,omitempty"` // 商品コード
	Price       *PriceInfo `json:"price" gorm:"-"`         // 現在の表示価格情報 (Nullable)
	// PrimaryImage     *ImageInfo      `json:"primary_image" gorm:"-"`
	Attributes []AttributeInfo `json:"attributes"gorm:"-"` // 対象SKUの属性リスト
	// AddedAt          time.Time       `json:"added_at"`           // お気に入り追加日時 (time.Time型)
	AddedAtFormatted string `json:"added_at_formatted"` // 通过Format格式化的日期

}

// --- 以下のDTOは他のAPIと共通化可能 ---

// PriceInfo 価格情報 (計算用数値と表示用文字列を含む)
// type PriceInfo struct {
// 	Amount                  float64  `json:"amount"`
// 	FormattedAmount         string   `json:"formatted_amount"`
// 	Type                    string   `json:"type"`
// 	TypeName                string   `json:"type_name"`
// 	OriginalAmount          *float64 `json:"original_amount,omitempty"`
// 	FormattedOriginalAmount *string  `json:"formatted_original_amount,omitempty"`
// }

// // ImageInfo 画像情報
// type ImageInfo struct {
// 	ID      int     `json:"id"`
// 	URL     string  `json:"url"`
// 	AltText *string `json:"alt_text,omitempty"`
// }

// // AttributeInfo 対象SKUの属性情報
// type AttributeInfo struct {
// 	AttributeID   int     `json:"attribute_id"`
// 	AttributeName string  `json:"attribute_name"`
// 	OptionID      *int    `json:"option_id,omitempty"`
// 	OptionValue   *string `json:"option_value,omitempty"`
// 	ValueString   *string `json:"value_string,omitempty"`
// 	SkuID         string  `json:"sku_id"` // 加入 SkuID 字段，用于与 SKU 对应
// 	// ... 他の value_xxx 型
// }

// // PaginationInfo ページネーション情報
type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalCount  int `json:"total_count"` // お気に入りの総件数
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
