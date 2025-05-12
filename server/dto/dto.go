package dto

// ProductInfoResponse APIのルートレスポンス構造
type ProductInfoResponse struct {
	ID                 string               `json:"id"`
	ProductCode        string               `json:"product_code,omitempty"`
	Name               string               `json:"name"`
	Description        *string              `json:"description,omitempty"`
	Status             string               `json:"status"`
	IsTaxable          bool                 `json:"is_taxable"`
	MetaTitle          *string              `json:"meta_title,omitempty"`
	MetaDescription    *string              `json:"meta_description,omitempty"`
	CreatedAtFormatted string               `json:"created_at_formatted"`
	UpdatedAtFormatted string               `json:"updated_at_formatted"`
	TargetSKUInfo      *TargetSKUInfo       `json:"target_sku_info,omitempty"` // 対象SKU情報 (Nullable)
	VariantOptions     []VariantOptionGroup `json:"variant_options,omitempty"` // 全バリエーション選択肢
	Category           *CategoryInfo        `json:"category,omitempty"`
	Brand              *BrandInfo           `json:"brand,omitempty"`
}

//
// TargetSKUInfo 対象SKUの詳細情報
type TargetSKUInfo struct {
	SkuID string `json:"sku_id"`

	Price *PriceInfo `json:"price,omitempty"`
	// PriceStr string     `json:"sku_price"`

	PrimaryImage *ImageInfo      `json:"primary_image,omitempty" gorm:"-"`
	Attributes   []AttributeInfo `json:"attributes"`
}

type PriceInfo struct {
	Amount                  float64  `json:"amount"`
	FormattedAmount         string   `json:"formatted_amount"`
	Type                    string   `json:"type,omitempty"`
	TypeName                string   `json:"type_name"`
	OriginalAmount          *float64 `json:"original_amount"`
	FormattedOriginalAmount *string  `json:"formatted_original_amount"`
}

// StockInfo 在庫状況概要
type StockInfo struct {
	Status           string  `json:"status"`
	StatusText       string  `json:"status_text"`
	DeliveryEstimate *string `json:"delivery_estimate,omitempty"`
}

// ImageInfo 画像情報
type ImageInfo struct {
	ID      int     `json:"id"`
	URL     string  `json:"url"`
	AltText *string `gorm:"column:alt_text" json:"alt_text"`
}

// AttributeInfo 対象SKUの属性情報
type AttributeInfo struct {
	AttributeID   int    `json:"attribute_id"`
	AttributeName string `json:"attribute_name"`
	// Value         *string `json:"value,omitempty"`
	// AttributeID   int     `json:"attribute_id"`
	// AttributeName string  `json:"attribute_name"`
	OptionID    *int    `json:"option_id,omitempty"`
	OptionValue *string `json:"option_value,omitempty"`
	ValueString *string `json:"value_string,omitempty"`
	// DisplayType   string          `json:"display_type"`
	SkuID string `json:"-" gorm:"-"`
	// ... 他の value_xxx 型
}

// VariantOptionGroup SKUバリエーション軸ごとの選択肢グループ
type VariantOptionGroup struct {
	AttributeID   int    `json:"attribute_id"`
	AttributeName string `json:"attribute_name"`
	AttributeCode string `json:"attribute_code"`
	// DisplayType   string          `json:"display_type"` // 'image', 'text', etc.
	Options []VariantOption `json:"options"`
}

// VariantOption 個々のバリエーション選択肢
type VariantOption struct {
	OptionID    int    `json:"option_id"`
	OptionValue string `json:"option_value"`
	OptionCode  string `json:"option_code"`
	// ImageURL     *string  `json:"image_url,omitempty"`      // Swatch画像等
	// IsSelectable bool     `json:"is_selectable"`            // 在庫等による選択可否
	LinkedSkuIDs []string `json:"linked_sku_ids,omitempty"` // 関連SKU (任意)
}
type ProductInfo struct {
	ProductID       string  `json:"id"`
	ProductCode     string  `json:"product_code,omitempty"`
	Name            string  `json:"name"`
	Description     string  `json:"-"`
	DescriptionStr  string  `json:"description"`
	IsTaxable       bool    `json:"is_taxable"`
	MetaTitle       *string `json:"meta_title"`
	MetaDescription *string `json:"meta_description"`
	SkuPrice        float64 `json:"-"`
	PriceStr        string  `json:"-"`
	// CategoryName    string    `json:"category_name"`
	// // AttributeIDs    int                `json:"attribute_ids"`
	TargetSKUInfo *TargetSKUInfo `json:"target_sku_info,omitempty" gorm:"-"`

	//临时字段，用于嵌套
	SkuID     string `json:"-"`
	PriceName string `json:"-"`
	PriceCode string `json:"-"`
	ImageURL  string `json:"-"`
	ImageAlt  string `json:"alt_text"`
	// AttributeIDs []int  `json:"attribute_ids"` // 用切片存储多个属性ID
}

// ProductVariantResponse 返回商品信息+变体选项
type ProductVariantResponse struct {
	ProductInfo ProductInfo          `json:"product_info"`
	Variants    []VariantOptionGroup `json:"variants"`
}

// CategoryInfo カテゴリ情報
type CategoryInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	ParentID *int   `json:"parent_id,omitempty"`
	// Breadcrumbs []Breadcrumb `json:"breadcrumbs,omitempty"`
}

// BrandInfo ブランド情報
type BrandInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse エラーレスポンス構造
type ProductErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
type ProductErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}
