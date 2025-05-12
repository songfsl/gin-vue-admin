package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ProductSkus struct {
	ProductID       string             `json:"id"`
	ProductCode     string             `json:"product_code,omitempty"`
	Name            string             `json:"name"`
	Description     string             `json:"-"`
	DescriptionStr  string             `json:"description"`
	IsTaxable       bool               `json:"is_taxable"`
	MetaTitle       *string            `json:"meta_title"`
	MetaDescription *string            `json:"meta_description"`
	SkuPrice        float64            `json:"-"`
	PriceStr        string             `json:"-"`
	TargetSKUInfo   *dto.TargetSKUInfo `json:"target_sku_info,omitempty" gorm:"-"`
	SkuID           string             `json:"-"`
	PriceName       string             `json:"-"`
	PriceCode       string             `json:"-"`
	ImageURL        string             `json:"-"`
	ImageAlt        string             `json:"-"`
}

type ProductSkusService struct{}

var ProductSkusApp = new(ProductSkusService)
var (
	ErrProductNotFound  = errors.New("product not found")
	ErrAlreadyFavorited = errors.New("product already favorited")
	ErrFavoriteNotFound = errors.New("favorite not found")
)

func TruncateString(s string, max int) string {
	if len([]rune(s)) > max {
		return string([]rune(s)[:max]) + "…"
	}
	return s
}

// var ProductSkusApp = new(ProductSkus)
func saveProductRedis(skuID string, productInfoResponse dto.ProductVariantResponse) {
	// 保存到redis

	//要先把结构体json化，才能存入redis
	redisData, err := json.Marshal(productInfoResponse)
	if err != nil {
		fmt.Println("err:", err)
	}
	//判断一下连没连上redis
	if global.GVA_REDIS == nil {
		fmt.Println("redis没初始化。")
	} else {
		fmt.Println("---------------------------------")
		fmt.Println("可以连接到redis, global.GVA_REDIS is:", global.GVA_REDIS)
		fmt.Println("")
	}
	ctx := context.Background()
	err = global.GVA_REDIS.Set(ctx, fmt.Sprintf("skuid:%s", skuID), redisData, 30*time.Minute).Err()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("---------------------------------")
	fmt.Println("已把skuid:", skuID, "的数据存入redis。")
	fmt.Println("")
}
func (P *ProductSkusService) CreateShippingAddress(UserId uint, req dto.ShippingAddressInput) error {
	//请求体,不是单独query
	db := global.GVA_DB
	//如果设置为默认地址，则要把原有的默认地址取消掉！！
	if req.IsDefault {
		err := db.Table("user_shipping_addresses").Where("user_id = ?", UserId).
			Update("is_default", false).Error
		if err != nil {
			return err
		}
	}
	//插入新地址
	type shipping struct {
		UserId        uint    `json:"user_id"`
		PostalCode    string  `json:"postal_code" binding:"required,max=10"`
		Prefecture    string  `json:"prefecture" binding:"required,max=50"`
		City          string  `json:"city" binding:"required,max=100"`
		AddressLine1  string  `json:"address_line1" binding:"required,max=255"`
		AddressLine2  *string `json:"address_line2,omitempty" binding:"max=255"`
		RecipientName string  `json:"recipient_name" binding:"required,max=100"`
		PhoneNumber   string  `json:"phone_number" binding:"required,max=20"`
		IsDefault     bool    `json:"is_default"`
	}
	newAddress := shipping{
		UserId:        UserId,
		PostalCode:    req.PostalCode,
		Prefecture:    req.Prefecture,
		City:          req.City,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		RecipientName: req.RecipientName,
		PhoneNumber:   req.PhoneNumber,
		IsDefault:     req.IsDefault,
	}
	err := db.Table("user_shipping_addresses").Create(&newAddress).Error
	return err

}
func (P *ProductSkusService) GetShippingAddress(UserId uint) (res dto.ShippingAddressListResponse, err error) {
	db := global.GVA_DB

	var shippingAddress []struct {
		AddressID     uint64  `json:"address_id"`
		PostalCode    string  `json:"postal_code"`
		Prefecture    string  `json:"prefecture"`
		City          string  `json:"city"`
		AddressLine1  string  `json:"address_line1"`
		AddressLine2  *string `json:"address_line2,omitempty"`
		RecipientName string  `json:"recipient_name"`
		PhoneNumber   string  `json:"phone_number"`
		IsDefault     bool    `json:"is_default"`
	}
	err = db.Table("user_shipping_addresses").
		Select("id AS address_id,postal_code,prefecture,city,address_line1,address_line2,recipient_name,phone_number,is_default,updated_at").
		Where("user_id = ?", UserId).
		Order("is_default DESC, updated_at DESC").
		Scan(&shippingAddress).Error
	if err != nil {
		return res, err
	}
	var results []dto.ShippingAddressInfo
	for _, v := range shippingAddress {
		result := dto.ShippingAddressInfo{
			AddressID:     v.AddressID,
			PostalCode:    v.PostalCode,
			Prefecture:    v.Prefecture,
			City:          v.City,
			AddressLine1:  v.AddressLine1,
			AddressLine2:  v.AddressLine2,
			RecipientName: v.RecipientName,
			PhoneNumber:   v.PhoneNumber,
			IsDefault:     v.IsDefault,
		}
		results = append(results, result)
	}
	res = dto.ShippingAddressListResponse{
		Addresses: results,
	}
	return res, nil

}
func (P *ProductSkusService) ChangeShippingAddress(UserId uint, AddressId int, req dto.ShippingAddressInput) error {
	db := global.GVA_DB
	//也是请求体
	// 确认该地址属于当前用户
	var count int64
	err := db.Table("user_shipping_addresses").
		Where("id = ? AND user_id = ?", AddressId, UserId).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("该地址不存在或不属于当前用户")
	}
	//如果设置为默认地址，则要把原有的默认地址取消掉！！
	if req.IsDefault {
		err := db.Table("user_shipping_addresses").Where("user_id = ?", UserId).
			Update("is_default", false).Error
		if err != nil {
			return err
		}
	}
	//更改地址
	type shipping struct {
		UserId        uint    `json:"user_id"`
		PostalCode    string  `json:"postal_code" binding:"required,max=10"`
		Prefecture    string  `json:"prefecture" binding:"required,max=50"`
		City          string  `json:"city" binding:"required,max=100"`
		AddressLine1  string  `json:"address_line1" binding:"required,max=255"`
		AddressLine2  *string `json:"address_line2,omitempty" binding:"max=255"`
		RecipientName string  `json:"recipient_name" binding:"required,max=100"`
		PhoneNumber   string  `json:"phone_number" binding:"required,max=20"`
		IsDefault     bool    `json:"is_default"`
	}
	newAddress := shipping{
		UserId:        UserId,
		PostalCode:    req.PostalCode,
		Prefecture:    req.Prefecture,
		City:          req.City,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		RecipientName: req.RecipientName,
		PhoneNumber:   req.PhoneNumber,
		IsDefault:     req.IsDefault,
	}
	//这里是updates，update更新单个字段，需要两个参数，updates一次可更新多个字段（结构体/map）
	err = db.Table("user_shipping_addresses").
		Where("id = ? AND user_id = ?", AddressId, UserId).Updates(newAddress).Error
	if err != nil {
		return err
	}

	return nil
}

func (P *ProductSkusService) DeleteShippingAddress(UserId uint, AddressId int) error {
	db := global.GVA_DB

	// 确认该地址属于当前用户
	var count int64
	err := db.Table("user_shipping_addresses").
		Where("id = ? AND user_id = ?", AddressId, UserId).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("该地址不存在或不属于当前用户")
	}

	// 删除该地址
	if err := db.Table("user_shipping_addresses").
		Where("id = ? ", AddressId).
		Delete(nil).Error; err != nil {
		return err
	}

	return nil
}

func (P *ProductSkusService) ChangeItemsInCart(UserId uint, SkuId string, Quantity int) error {
	db := global.GVA_DB
	// 先检查 SKU 是否存在
	var skuCount int64
	if err := db.Table("product_skus").
		Where("id = ?", SkuId).
		Count(&skuCount).Error; err != nil {
		return err
	}
	if skuCount == 0 {
		return ErrProductNotFound
	}
	//查询cart里sku
	type CartItem struct {
		UserId   uint
		SkuId    string
		Quantity int
	}
	var existing CartItem
	err := db.Table("user_cart_items").Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Take(&existing).Error
	if err != nil {
		return ErrProductNotFound
	}
	//查询本身库存信息(不加购物车之前)
	var tempQuantity struct {
		Quantity         int `gorm:"column:quantity"`
		ReservedQuantity int `gorm:"column:reserved_quantity"`
	}
	err = db.Table("inventory").
		Select("quantity,reserved_Quantity").
		Where("sku_id = ?", SkuId).Scan(&tempQuantity).Error
	if err != nil {
		return err
	}
	//加入购物车，判断加入数量不能超过available!!
	available := tempQuantity.Quantity - tempQuantity.ReservedQuantity
	if Quantity > available {
		return fmt.Errorf("超过可购买数量")
	}

	err = db.Table("user_cart_items").Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Update("quantity", Quantity).Error

	return err
}

func (P *ProductSkusService) DeleteItemsFromCart(UserId uint, SkuId string) error {
	db := global.GVA_DB

	// 先检查 SKU 是否存在
	var skuCount int64
	if err := db.Table("product_skus").
		Where("id = ?", SkuId).
		Count(&skuCount).Error; err != nil {
		return err
	}
	if skuCount == 0 {
		return ErrProductNotFound
	}
	// 检查是否已加入
	var ItemCount int64
	if err := db.Table("user_cart_items").
		Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Count(&ItemCount).Error; err != nil {
		return err
	}
	if ItemCount == 0 {
		return ErrFavoriteNotFound
	}
	// 执行删除
	if err := db.Table("user_cart_items").
		Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Delete(nil).Error; err != nil {
		return err
	}
	return nil
}
func (P *ProductSkusService) AddItemsIntoCart(UserId uint, SkuId string, Quantity int) error {
	// 查询 product_skus 表，看 sku_id 是否存在
	//建议用count
	var count int64
	db := global.GVA_DB
	err := db.Table("product_skus").
		Select("product_skus.id").
		Where("id = ?", SkuId).
		Count(&count).Error
	if err != nil {
		return err
	} else if count == 0 {
		return ErrProductNotFound
	}
	//查询本身库存信息(不加购物车之前)
	var tempQuantity struct {
		Quantity         int `gorm:"column:quantity"`
		ReservedQuantity int `gorm:"column:reserved_quantity"`
	}
	err = db.Table("inventory").
		Select("quantity,reserved_Quantity").
		Where("sku_id = ?", SkuId).Scan(&tempQuantity).Error
	if err != nil {
		return err
	}
	//加入购物车，判断加入数量不能超过available!!
	available := tempQuantity.Quantity - tempQuantity.ReservedQuantity
	if available <= 0 {
		return fmt.Errorf("ErrInsufficientStock")
	}
	type CartItem struct {
		UserId   uint
		SkuId    string
		Quantity int
	}
	//查询cart里是否已存在，存在就增加数量（更新）
	var existing CartItem
	err = db.Table("user_cart_items").Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Take(&existing).Error
	if err == nil {
		//如果没有报错，说明存在，就更新数量update
		//如果有错误，说明不存在直接create
		new := existing.Quantity + Quantity
		if new > available {
			return fmt.Errorf("超过可购买数量")
		} //加入表
		err = db.Table("user_cart_items").Where("user_id = ? AND sku_id = ?", UserId, SkuId).
			Update("quantity", new).Error
	} else if err != nil {
		if Quantity > available {
			return fmt.Errorf("超过可购买数量")
		} // 3. 再插入记录，指定表名！(匿名结构体插入数据gorm无法判断表名)
		err = db.Table("user_cart_items").Create(&CartItem{
			UserId:   UserId,
			SkuId:    SkuId,
			Quantity: Quantity,
		}).Error

	}

	return err
}
func (P *ProductSkusService) GetCartItems(UserId uint) (res dto.CartResponse, err error) {
	//去重问题待解决。。。
	db := global.GVA_DB
	var totalCount int64
	// 计算用户购物车中的商品数量总和
	err = db.Table("user_cart_items").
		Select("SUM(user_cart_items.quantity) AS total_quantity").
		Joins("JOIN product_skus ON product_skus.id = user_cart_items.sku_id").
		Where("user_cart_items.user_id = ?", UserId).
		Scan(&totalCount).Error
	if err != nil {
		return res, err
	}
	var cartItems []struct {
		SkuId                   string
		ProductID               string
		ProductName             string
		ProductCode             string
		Quantity                int
		Amount                  float64
		FormattedAmount         string
		Type                    string
		TypeName                string
		OriginalAmount          *float64
		FormattedOriginalAmount *string
		SubtotalFormatted       string
		ID                      int
		URL                     string
		AltText                 *string
		AttributeID             int
		AttributeName           string
		OptionID                *int
		OptionValue             *string
		ValueString             *string
		StockStatus             string
	}
	db.Table("product_skus").
		Select(`
		product_skus.id AS sku_id,
		products.id AS product_id,
		products.name AS product_name,
		products.product_code,
		user_cart_items.quantity,
		prices.price AS amount,
		price_types.type_code AS type,
		price_types.name AS name,
		(prices.price * user_cart_items.quantity) AS subtotal_formatted,
		sku_images.thumbnail_url AS url,
		sku_images.alt_text AS alt_text,
		sku_images.id AS id,
		user_cart_items.updated_at,
		MAX(attributes.id) AS attribute_id,
		MAX(attributes.name) AS attribute_name,
		CASE
			WHEN attribute_options.value IS NOT NULL THEN attribute_options.value
			WHEN sku_values.value_string IS NOT NULL THEN sku_values.value_string
			WHEN sku_values.value_number IS NOT NULL THEN CAST(sku_values.value_number AS CHAR)
			WHEN sku_values.value_boolean IS NOT NULL THEN CAST(sku_values.value_boolean AS CHAR)
			ELSE ''
		END AS value_string,
		MAX(attribute_options.id) AS option_id,
		MAX(attribute_options.value) AS option_value,
		CASE 
			WHEN inventory.quantity = 0 THEN 'out_of_stock'
			WHEN inventory.quantity < 10 THEN 'low_stock'
			ELSE 'available'
		END AS stock_status
	`).
		Joins("JOIN user_cart_items ON product_skus.id = user_cart_items.sku_id").
		Joins("LEFT JOIN products ON product_skus.product_id = products.id").
		Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Joins("LEFT JOIN sku_values ON product_skus.id = sku_values.sku_id").
		Joins("JOIN attributes ON sku_values.attribute_id = attributes.id").
		Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
		Joins("LEFT JOIN sku_images ON product_skus.id = sku_images.sku_id").
		Joins("LEFT JOIN inventory ON product_skus.id = inventory.sku_id").
		Where("user_cart_items.user_id = ?", UserId).
		Group(`
		product_skus.id,
		products.id,
		products.name,
		products.product_code,
		user_cart_items.quantity,
		prices.price,
		price_types.type_code,
		price_types.name,
		sku_images.thumbnail_url,
		sku_images.alt_text,
		sku_images.id,
		user_cart_items.updated_at,
		attributes.id,
		inventory.quantity
	`).
		Scan(&cartItems)

	if err != nil {
		return res, err
	}
	var results []dto.CartItemInfo
	var totalAmount float64
	var totalAmountFormatted string
	//map聚合
	cartMap := make(map[string]*dto.CartItemInfo)
	//price需要格式化
	for _, p := range cartItems {
		//如果不存在这个skuid，就创建一个
		if _, ok := cartMap[p.SkuId]; !ok {
			formatted := fmt.Sprintf("%s円", humanize.Commaf(float64(p.Amount)))
			priceFloat, _ := strconv.ParseFloat(p.SubtotalFormatted, 64) // 注意处理错误
			subtotalFormatted := fmt.Sprintf("%s円", humanize.Commaf(priceFloat))
			price := p.Amount // float64
			quantity := p.Quantity
			totalAmount += price * float64(quantity)

			cartMap[p.SkuId] = &dto.CartItemInfo{
				SkuID:       p.SkuId,
				ProductID:   p.ProductID,
				ProductCode: p.ProductCode,
				ProductName: p.ProductName,
				Quantity:    p.Quantity,
				Price: &dto.PriceInfo{
					Amount:                  p.Amount,
					FormattedAmount:         formatted,
					Type:                    p.Type,
					TypeName:                p.TypeName,
					OriginalAmount:          p.OriginalAmount,
					FormattedOriginalAmount: p.FormattedOriginalAmount,
				},
				SubtotalFormatted: subtotalFormatted,
				PrimaryImage: &dto.ImageInfo{
					ID:      p.ID,
					URL:     p.URL,
					AltText: p.AltText,
				},
				Attributes:  []dto.AttributeInfo{},
				StockStatus: p.StockStatus,
			}
		}
		//加属性
		if p.AttributeID != 0 {
			attr := dto.AttributeInfo{
				AttributeID:   p.AttributeID,
				AttributeName: p.AttributeName,
				OptionID:      p.OptionID,
				OptionValue:   p.OptionValue,
				ValueString:   p.ValueString,
			}
			cartMap[p.SkuId].Attributes = append(cartMap[p.SkuId].Attributes, attr)
		}
	}
	//map转【】
	for _, item := range cartMap {
		results = append(results, *item)
	}
	totalAmountFormatted = fmt.Sprintf("%s円", humanize.Commaf(totalAmount))

	// 构建最终响应 DTO
	res = dto.CartResponse{
		Items:                results,
		TotalItemsCount:      int(totalCount),
		TotalAmount:          totalAmount,
		TotalAmountFormatted: totalAmountFormatted, // 格式化金额函数
	}
	return res, err

}
func (P *ProductSkusService) GetViewedHistory(UserId uint, Page int, Limit int) (res dto.ViewedSKUListResponse, err error) {
	db := global.GVA_DB

	// 获取总数
	var totalCount int64
	err = db.Table("user_viewed_skus").
		Where("user_id = ?", UserId).
		Count(&totalCount).Error
	if err != nil {
		return res, err
	}
	var viewedSkuInfo []struct {
		SkuId             string
		ProductID         string
		ProductName       string
		ProductCode       string
		MinPrice          float64
		MaxPrice          float64
		ID                int
		URL               string
		AltText           *string
		AverageRating     *float64
		ReviewCount       *int
		ViewedAtFormatted string
	}
	offset := (Page - 1) * Limit
	err = db.Table("product_skus").
		Select(`
		product_skus.id AS sku_id,
		products.id as product_id,
		products.name as product_name,
		products.product_code,
		MAX(prices.price) as max_price,
		  MIN(prices.price) as min_price,
		  sku_images.sku_id,
		  sku_images.thumbnail_url as url,
		sku_images.alt_text,
		sku_images.id,
		review_summaries.average_rating,
		review_summaries.review_count,
		DATE_FORMAT(user_viewed_skus.viewed_at, '%Y年%m月%d日 %H:%i:%s') AS viewed_at_formatted
	`).
		Joins("JOIN user_viewed_skus ON product_skus.id = user_viewed_skus.sku_id").
		Joins("LEFT JOIN products ON product_skus.product_id = products.id ").
		Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
		Joins("LEFT JOIN review_summaries ON products.id = review_summaries.product_id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Joins(`
LEFT JOIN (
 SELECT
      sku_id,
      MIN(thumbnail_url) AS thumbnail_url,
      MIN(alt_text) AS alt_text,
      MAX(id) AS id
    FROM sku_images
    GROUP BY sku_id
) AS sku_images ON product_skus.id= sku_images.sku_id
`).
		Where(" user_viewed_skus.user_id = ?", UserId).
		Order("user_viewed_skus.viewed_at DESC").
		Group(`
		product_skus.id,
	products.id,
	products.product_code,
	products.name,
	review_summaries.average_rating,
	review_summaries.review_count,
	user_viewed_skus.viewed_at
`).
		Limit(Limit).
		Offset(offset).
		Scan(&viewedSkuInfo).Error
	if err != nil {
		return res, err
	}
	var results []dto.ViewedSKUInfo
	//price需要格式化
	for _, p := range viewedSkuInfo {
		formatted := fmt.Sprintf("%s円", humanize.Commaf(float64(p.MinPrice)))
		//有min和max
		if p.MinPrice != p.MaxPrice {
			formatted = fmt.Sprintf("%s円 ~ %s円", humanize.Commaf(float64(p.MinPrice)), humanize.Commaf(float64(p.MaxPrice)))
		}
		result := dto.ViewedSKUInfo{
			SkuID:               p.SkuId,
			ProductID:           p.ProductID,
			ProductCode:         p.ProductCode,
			ProductName:         p.ProductName,
			PriceRangeFormatted: formatted,
			PrimaryImage: &dto.ImageInfo{
				ID:      p.ID,
				URL:     p.URL,
				AltText: p.AltText,
			},
			ReviewSummary: &dto.ReviewSummaryInfo{
				AverageRating: *p.AverageRating,
				ReviewCount:   *p.ReviewCount,
			},
			ViewedAtFormatted: p.ViewedAtFormatted,
		}
		results = append(results, result)
	}
	// 分页信息
	totalPages := int((totalCount + int64(Limit) - 1) / int64(Limit))
	pagination := dto.PaginationInfo{
		CurrentPage: Page,
		Limit:       Limit,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
	}

	res = dto.ViewedSKUListResponse{
		History:    results,
		Pagination: pagination,
	}
	return res, nil
}

func (P *ProductSkusService) AddViewedSkus(UserId uint, SkuId string) error {
	// 1. 查询 product_skus 表，看 sku_id 是否存在
	//建议用count
	var count int64
	db := global.GVA_DB
	err := db.Table("product_skus").
		Select("product_skus.id").
		Where("id = ?", SkuId).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		// return res, errors.New("ProductCode 不能为空")
		return ErrProductNotFound
	}
	// 2. 查询 user_viewed_skus 表是否已存在记录
	var existing int64
	err = db.Table("user_viewed_skus").
		Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Count(&existing).Error
	if err != nil {
		return err
	}
	if existing > 0 {
		return ErrAlreadyFavorited
	}

	// 3. 插入浏览记录，指定表名！(匿名结构体插入数据gorm无法判断表名)
	err = db.Table("user_viewed_skus").Create(&struct {
		UserId uint
		SkuId  string
	}{
		UserId: UserId,
		SkuId:  SkuId,
	}).Error

	return err
}

func (P *ProductSkusService) GetStaffCoordinate(ProductCode string, Limit int) (res dto.CoordinateSetTeaserListResponse, err error) {
	var coordinateSetTeaser []struct {
		SetID                string
		SetThemeImageURL     *string
		ContributorNickname  string
		ContributorAvatarURL *string
		ContributorStoreName *string
	}
	db := global.GVA_DB
	err = db.Table("coordinate_sets").
		Select(`
            coordinate_sets.id as set_id,
			coordinate_sets.theme_image_url as set_theme_image_url,
             coordinate_sets.contributor_nickname,
			 coordinate_sets.contributor_avatar_url,
			 coordinate_sets.contributor_store_name,
			  coordinate_set_items.product_id,
			  coordinate_sets.sort_order,
			  products.id,
			  products.product_code,
			   coordinate_set_items.coordinate_set_id
		`).
		Joins("LEFT JOIN coordinate_set_items ON  coordinate_sets.id =  coordinate_set_items.coordinate_set_id ").
		Joins("LEFT JOIN products ON  coordinate_set_items.product_id = products.id").
		Where("products.product_code = ?", ProductCode).
		Order(" coordinate_sets.sort_order ASC ").
		Limit(Limit).
		Scan(&coordinateSetTeaser).Error
	if err != nil {
		return res, err
	}

	var results []dto.CoordinateSetTeaserInfo
	for _, v := range coordinateSetTeaser {
		result := dto.CoordinateSetTeaserInfo{
			SetID:                v.SetID,
			SetThemeImageURL:     v.SetThemeImageURL,
			ContributorNickname:  v.ContributorNickname,
			ContributorAvatarURL: v.ContributorAvatarURL,
			ContributorStoreName: v.ContributorStoreName,
		}
		results = append(results, result)
		//放到response
		res = dto.CoordinateSetTeaserListResponse{
			Coordinates: results,
		}
	}
	return res, nil

}

func (P *ProductSkusService) GetRelatedProduct(ProductCode string, Limit int) (res dto.RelatedProductListResponse, err error) {
	var relatedProduct []struct {
		ProductID         string
		ProductCode       string
		ProductName       string
		MinPrice          float64
		MaxPrice          float64
		IsOnSale          bool
		AverageRating     *float64
		ReviewCount       *int
		ThumbnailImageURL *string

		CategoryID       uint
		CategoryName     string
		CategoryLevel    int
		CategoryParentID int
		CategoryLinks    *string
	}

	db := global.GVA_DB
	//思路分析：1.通过productcode取categoryid
	var categoryId uint
	err = db.Table("products").
		Select("category_id").Where("product_code = ?", ProductCode).Scan(&categoryId).Error
	if err != nil {
		return res, err
	}
	//2.用categoryid取info
	//product下有多个sku，价格范围
	//注意是要一个图片，要不然会重复（不知道需要需要改..
	err = db.Table("products").
		Select(`
			products.id as product_id,
			products.product_code,
			products.name as product_name,
			review_summaries.average_rating,
			review_summaries.review_count,
			sku_images.thumbnail_url as thumbnail_image_url,
			MAX(CASE WHEN prices.price_type_id = 2 THEN 1 ELSE 0 END) as is_on_sale,
			 MAX(prices.price) as max_price,
			  MIN(prices.price) as min_price,
			  categories.id as category_id,
			 categories.name as category_name,
			 categories.parent_id as category_parent_id,
			 categories.level as category_level,
			  (
        SELECT GROUP_CONCAT(c.name ORDER BY c.id)
        FROM categories c
        WHERE c.parent_id = categories.parent_id AND c.id != categories.id
    ) AS category_links
		`).
		Joins("LEFT JOIN product_skus ON products.id = product_skus.product_id ").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
		Joins("LEFT JOIN review_summaries ON products.id = review_summaries.product_id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Joins(`
  LEFT JOIN (
    SELECT sku_id, MIN(thumbnail_url) AS thumbnail_url
    FROM sku_images
    GROUP BY sku_id
  ) AS sku_images ON products.default_sku_id = sku_images.sku_id
`). // Joins("LEFT JOIN sku_images ON products.default_sku_id = sku_images.sku_id ").
		Where("products.category_id = ?", categoryId).
		Group(`
		products.id,
		products.product_code,
		products.name,
		review_summaries.average_rating,
		review_summaries.review_count,
		sku_images.thumbnail_url,
		 categories.id, categories.name, categories.parent_id, categories.level
	`).
		Limit(Limit).
		// Having("price_range_formatted IS NOT NULL").
		Scan(&relatedProduct).Error
	if err != nil {
		return res, err
	}
	//抽关联category(先单独抽，抽出来再尝试优化和合到上面)
	// 拆分类信息（只取第一个即可）
	first := relatedProduct[0]
	category := dto.RelatedCategoryInfo{
		CategoryID:       first.CategoryID,
		CategoryName:     first.CategoryName,
		CategoryLevel:    first.CategoryLevel,
		CategoryParentID: first.CategoryParentID,
	}
	// 转换 *string 到 []string
	if first.CategoryLinks != nil {
		category.CategoryLinks = strings.Split(*first.CategoryLinks, ",")
	}
	var results []dto.RelatedProductInfo
	//price需要格式化
	for _, p := range relatedProduct {
		formatted := fmt.Sprintf("%s円", humanize.Commaf(float64(p.MinPrice)))
		//有min和max
		if p.MinPrice != p.MaxPrice {
			formatted = fmt.Sprintf("%s円 ~ %s円", humanize.Commaf(float64(p.MinPrice)), humanize.Commaf(float64(p.MaxPrice)))
		}

		//塞infodto
		result := dto.RelatedProductInfo{
			ProductID:           p.ProductID,
			ProductCode:         p.ProductCode,
			ProductName:         p.ProductName,
			PriceRangeFormatted: formatted,
			IsOnSale:            p.IsOnSale,
			ReviewSummary: &dto.ReviewSummaryInfo{
				AverageRating: *p.AverageRating,
				ReviewCount:   *p.ReviewCount,
			},
			ThumbnailImageURL: p.ThumbnailImageURL,
			//linkedSkuIDs := strings.Split(v.LinkedSkuIDs, ",")
		}
		results = append(results, result)
		//放到response
		res = dto.RelatedProductListResponse{
			RelatedProducts: results,
			RelatedCategory: category,
		}
	}

	return res, nil

}

func (P *ProductSkusService) GetTargetProductSkus(SkuID string, ProductID string) (results ProductSkus, err error) {

	db := global.GVA_DB

	//思路分析：
	// 清理输入,空格总是报错，chat叫我加这个
	SkuID = strings.TrimSpace(SkuID)
	ProductID = strings.TrimSpace(ProductID)

	// 1.如果 SkuID 存在，直接查 SKU 信息
	if SkuID != "" {
		err = db.Table("product_skus").
			Select(`
			product_skus.product_id,
		product_skus.id AS sku_id, 
		product_skus.sku_code, 
		product_skus.status As status,
		products.name AS name, 
		products.product_code, 
		products.description, 
		products.is_taxable, 
		products.created_at,
		products.updated_at,
		products.meta_title, 
		products.meta_description, 
		categories.id AS category_id, 
		categories.name AS category_name, 
		categories.level AS category_level, 
		categories.parent_id AS category_parent_id,
		sku_images.id AS image_id,
		sku_images.image_url AS image_url,
		sku_images.alt_text AS alt_text,
		prices.id AS price_id,
		prices.price AS sku_price,
		prices.start_date AS price_start_date,
		prices.end_date AS  price_end_date,
		price_types.type_code AS price_type,
		price_types.name AS price_type_name`).
			Joins("LEFT JOIN products  ON product_skus.product_id = products.id").
			Joins("LEFT JOIN categories  ON products.category_id = categories.id").
			Joins("LEFT JOIN prices  ON prices.sku_id = product_skus.id").
			Joins("LEFT JOIN price_types  ON price_types.id = prices.price_type_id").
			Joins("LEFT JOIN sku_images ON sku_images.sku_id = product_skus.id ").
			Where("product_skus.id = ?", SkuID).
			Find(&results).Error
		// results.CreatedAtStr = results.CreatedAt.Format("2006-01-02 15:04:05")
		// results.UpdatedAtStr = results.UpdatedAt.Format("2006-01-02 15:04:05")
		results.PriceStr = humanize.Commaf(results.SkuPrice) + "円"
		// title := "这是一个非常非常长的标题，应该被截断"
		results.DescriptionStr = TruncateString(results.Description, 20)
		//部分嵌套
		var attributes []dto.AttributeInfo
		//value 有两个一个是string另一个是number，如果需要用sql处理非空，费劲
		err = db.Table("sku_values").
			Select(`attributes.id AS attribute_id,
				attributes.name AS attribute_name,
			  CASE
				WHEN attribute_options.value IS NOT NULL THEN attribute_options.value
				WHEN sku_values.value_string IS NOT NULL THEN sku_values.value_string
				WHEN sku_values.value_number IS NOT NULL THEN CAST(sku_values.value_number AS CHAR)
				WHEN sku_values.value_boolean IS NOT NULL THEN CAST(sku_values.value_boolean AS CHAR)
				ELSE ''
			END AS value`).
			Joins("JOIN attributes ON sku_values.attribute_id = attributes.id").
			Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
			Where("sku_values.sku_id = ?", SkuID).

			// 按照sort_order升序排序
			Order("attributes.sort_order ASC").
			Scan(&attributes).Error

		results.TargetSKUInfo = &dto.TargetSKUInfo{
			SkuID: results.SkuID,
			Price: &dto.PriceInfo{
				Amount:                  results.SkuPrice,
				FormattedAmount:         results.PriceStr,
				Type:                    results.PriceCode,
				TypeName:                results.PriceName,
				OriginalAmount:          nil,
				FormattedOriginalAmount: nil,
			},
			PrimaryImage: &dto.ImageInfo{
				URL:     results.ImageURL,
				AltText: &results.ImageAlt,
			},
			Attributes: attributes,
		}
		err = fmt.Errorf("输入有误，无法查询")
		return results, errors.New("invalid product id")
	}

	//2. 如果没有给 SkuID，但有 ProductID
	//分两步查的，一步写一直报错找不到在哪，两步清晰一些

	if ProductID != "" {
		var defaultSkuID string
		err = db.Table("products").
			Select("default_sku_id").
			Where("id = ?", ProductID).
			Scan(&defaultSkuID).Error

		if err != nil || defaultSkuID == "" {
			return results, errors.New("invalid product id")
		}

		err = db.Table("product_skus").
			Select(`
		product_skus.id AS sku_id, 
        product_skus.sku_code, 
        product_skus.product_id,
		product_skus.status As status,
        products.name AS name, 
        products.product_code, 
        products.description, 
        products.is_taxable, 
        products.created_at,
        products.updated_at,
        products.meta_title, 
        products.meta_description, 
        categories.id AS category_id, 
        categories.name AS category_name, 
        categories.level AS category_level, 
        categories.parent_id AS category_parent_id,
        sku_images.image_url AS image_url,
        sku_images.alt_text AS image_alt_text,
        price_types.id,
        prices.price AS sku_price,
		price_types.type_code as price_code,
		price_types.name as price_name,
        sku_images.image_url AS image_url,
        sku_images.alt_text AS image_alt_text
    `).
			Joins("LEFT JOIN products ON products.default_sku_id = product_skus.id").
			Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
			Joins("LEFT JOIN categories ON products.category_id = categories.id").
			Joins("LEFT JOIN sku_values ON product_skus.id = sku_values.sku_id").
			Joins("LEFT JOIN attributes ON sku_values.attribute_id = attributes.id").
			Joins("LEFT JOIN sku_images ON sku_images.sku_id = product_skus.id").
			Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
			Where("product_skus.id = ?", defaultSkuID).
			First(&results).Error

		// results.CreatedAtStr = results.CreatedAt.Format("2006-01-02 15:04:05")
		// results.UpdatedAtStr = results.UpdatedAt.Format("2006-01-02 15:04:05")
		results.DescriptionStr = TruncateString(results.Description, 20)
		results.PriceStr = humanize.Commaf(results.SkuPrice) + "円"
		//部分嵌套,attribute部分直接调用下面方法
		var attributes []dto.AttributeInfo
		//value 有两个一个是string另一个是number，如果需要用sql处理非空，费劲
		err = db.Table("sku_values").
			Select(`attributes.id AS attribute_id,
				attributes.name AS attribute_name,
			  CASE
				WHEN attribute_options.value IS NOT NULL THEN attribute_options.value
				WHEN sku_values.value_string IS NOT NULL THEN sku_values.value_string
				WHEN sku_values.value_number IS NOT NULL THEN CAST(sku_values.value_number AS CHAR)
				WHEN sku_values.value_boolean IS NOT NULL THEN CAST(sku_values.value_boolean AS CHAR)
				ELSE ''
			END AS value`).
			Joins("JOIN attributes ON sku_values.attribute_id = attributes.id").
			Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
			Where("sku_values.sku_id = ?", defaultSkuID).

			// 按照sort_order升序排序
			Order("attributes.sort_order ASC").
			Scan(&attributes).Error

		results.TargetSKUInfo = &dto.TargetSKUInfo{
			SkuID: results.SkuID,
			Price: &dto.PriceInfo{
				Amount:                  results.SkuPrice,
				FormattedAmount:         results.PriceStr,
				Type:                    results.PriceCode,
				TypeName:                results.PriceName,
				OriginalAmount:          nil,
				FormattedOriginalAmount: nil,
			},
			PrimaryImage: &dto.ImageInfo{
				URL:     results.ImageURL,
				AltText: &results.ImageAlt,
			},
			Attributes: attributes,
		}

		return results, err
	}

	//3.如果啥都没给，返回空结果和错误
	err = fmt.Errorf("SkuID 和 ProductID 都为空，无法查询")
	return results, err
}

// 思路分析：可以套用skuid/productid查询商品详情（对它的升级版）
// 1.查TargetSKUInfo
// 2.查VariantOptions
// 3.查Category
// 4.返回dto
func (P *ProductSkusService) GetVariantOptions(ProductID string) (res dto.ProductVariantResponse, err error) {
	db := global.GVA_DB

	// 去除 ProductID 的空格,否则报错
	ProductID = strings.TrimSpace(ProductID)
	//最终返回的结构体
	var productInfoResponse dto.ProductVariantResponse

	//通过productid，得到默认skuid
	var defaultSkuID string
	err = db.Table("products").
		Select("default_sku_id").
		Where("id = ?", ProductID).
		Scan(&defaultSkuID).Error

	if err != nil || defaultSkuID == "" {

		return res, errors.New("invalid product id")
	}
	// 根据skuid去redis里先查找，有数据直接返回，没有再去找数据库
	ctx := context.Background()
	data, err := global.GVA_REDIS.Get(ctx, fmt.Sprintf("skuid:%s", defaultSkuID)).Result()
	if err != nil {
		fmt.Print(err)
	} else if data != "" {
		err = json.Unmarshal([]byte(data), &productInfoResponse)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("---------------------------------")
		fmt.Println("已直接从redis获取skuid:", defaultSkuID, "的数据。")
		fmt.Println("")
		return productInfoResponse, err

	}

	//通过skuid显示商品详情

	var productInfo dto.ProductInfo
	// var results ProductSkus
	err = db.Table("product_skus").
		Select(`
				product_skus.id AS sku_id, 
        product_skus.sku_code, 
        product_skus.product_id,
		product_skus.status As status,
        products.name AS name, 
        products.product_code, 
        products.description, 
        products.is_taxable, 
        products.created_at,
        products.updated_at,
        products.meta_title, 
        products.meta_description, 
        categories.id AS category_id, 
        categories.name AS category_name, 
        categories.level AS category_level, 
        categories.parent_id AS category_parent_id,
        price_types.id,
        prices.price AS sku_price,
		price_types.type_code as price_code,
		price_types.name as price_name,
        sku_images.image_url AS image_url,
        sku_images.alt_text AS alt_text,
		sku_images.image_type
    `).
		Joins("LEFT JOIN products ON products.default_sku_id = product_skus.id").
		Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Joins("LEFT JOIN sku_values ON product_skus.id = sku_values.sku_id").
		Joins("LEFT JOIN attributes ON sku_values.attribute_id = attributes.id").
		Joins("LEFT JOIN sku_images ON sku_images.sku_id = product_skus.id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Where("product_skus.id = ?", defaultSkuID).
		Scan(&productInfo).Error

	// results.CreatedAtStr = results.CreatedAt.Format("2006-01-02 15:04:05")
	// results.UpdatedAtStr = results.UpdatedAt.Format("2006-01-02 15:04:05")
	productInfo.DescriptionStr = TruncateString(productInfo.Description, 20)
	productInfo.PriceStr = humanize.Commaf(productInfo.SkuPrice) + "円"
	//部分嵌套
	var attributes []dto.AttributeInfo
	//value 有两个一个是string另一个是number，如果需要用sql处理非空，费劲
	err = db.Table("sku_values").
		Select(`attributes.id AS attribute_id,
            attributes.name AS attribute_name,
          CASE
            WHEN attribute_options.value IS NOT NULL THEN attribute_options.value
            WHEN sku_values.value_string IS NOT NULL THEN sku_values.value_string
            WHEN sku_values.value_number IS NOT NULL THEN CAST(sku_values.value_number AS CHAR)
            WHEN sku_values.value_boolean IS NOT NULL THEN CAST(sku_values.value_boolean AS CHAR)
            ELSE ''
        END AS value`).
		Joins("JOIN attributes ON sku_values.attribute_id = attributes.id").
		Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
		Where("sku_values.sku_id = ?", defaultSkuID).

		// 按照sort_order升序排序
		Order("attributes.sort_order ASC").
		Scan(&attributes).Error

	var originalPrice float64
	err = db.Table("prices").
		Select("price").
		Where("sku_id = ? AND price_type_id = 1", defaultSkuID).
		Scan(&originalPrice).Error

	if err != nil {
		return res, err
	}

	// 查询是否有有效的促销价（price_type_id = 2）
	var salePriceInfo struct {
		Price    float64
		TypeCode string
		TypeName string
	}
	err = db.Table("prices").
		Select("prices.price, price_types.type_code, price_types.name").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Where("sku_id = ? AND price_type_id = 2 AND start_date <= ? AND end_date >= ?", defaultSkuID, time.Now(), time.Now()).
		Scan(&salePriceInfo).Error
	if err != nil {
		return res, err
	}

	// 构建 PriceInfo
	var priceInfo dto.PriceInfo
	if salePriceInfo.Price > 0 {
		// 使用
		priceInfo = dto.PriceInfo{
			Amount:                  salePriceInfo.Price,
			FormattedAmount:         humanize.Commaf(salePriceInfo.Price) + "円",
			Type:                    salePriceInfo.TypeCode,
			TypeName:                "セール価格",
			OriginalAmount:          nil,
			FormattedOriginalAmount: nil,
		}
	} else {
		// 使用原价
		priceInfo = dto.PriceInfo{
			Amount:                  originalPrice,
			FormattedAmount:         humanize.Commaf(originalPrice) + "円",
			Type:                    "original",
			TypeName:                "通常価格",
			OriginalAmount:          nil,
			FormattedOriginalAmount: nil,
		}
	}

	// 构建 TargetSKUInfo
	productInfo.TargetSKUInfo = &dto.TargetSKUInfo{
		SkuID: productInfo.SkuID,
		Price: &priceInfo,
		PrimaryImage: &dto.ImageInfo{
			URL:     productInfo.ImageURL,
			AltText: &productInfo.ImageAlt,
		},
		Attributes: attributes,
	}

	res.ProductInfo = productInfo

	//1.用productid查所有skuid
	var skuIds []string //因为一个productid可能对应很多个skuid，切片比较便于管理

	err = db.Table("product_skus").Select("id").Where("product_id =?", ProductID).
		Pluck("id", &skuIds).Error
	//pluck查询一列
	//判断有没有sku
	if err != nil {
		return res, err
	}
	//要用len
	if len(skuIds) == 0 {
		return res, err
	}
	// 2. 查询属性和对应的sku_id

	//2.查属性变体信息
	//思路分析：建立临时结构体,查sku属性信息，再按照attributeid进行分组
	//嵌套
	//定义临时结构体，接收查到的数据
	var variantOptions []struct {
		AttributeID   int    `json:"attribute_id"`
		AttributeName string `json:"attribute_name"`
		AttributeCode string `json:"attribute_code"`
		OptionID      int    `json:"option_id"`
		OptionValue   string `json:"option_value"`
		OptionCode    string `json:"option_code"`
		LinkedSkuIDs  string `json:"linked_sku_ids,omitempty"` // 関連SKU (任意)
	}

	err = db.Table("sku_values").
		Select(`
		attributes.id as attribute_id,
		attributes.name as attribute_name,
		attribute_options.value as option_value,
		attributes.attribute_code,
		attribute_options.option_code,
		attribute_options.id as option_id,
		category_attributes.attribute_id as category_attr_id,
		GROUP_CONCAT(DISTINCT sku_values.sku_id) as linked_sku_ids
	`).
		Joins("LEFT JOIN attributes ON sku_values.attribute_id = attributes.id").
		Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
		Joins("LEFT JOIN category_attributes ON attributes.id = category_attributes.attribute_id").
		Where("sku_values.sku_id IN ? AND category_attributes.is_variant_attribute = ?", skuIds, 1).
		Group("attributes.id, attribute_options.id").
		Order("attribute_options.attribute_id ASC").
		Scan(&variantOptions).Error
	//按属性以及选择分组

	if err != nil {
		return res, err
	}
	// sql分组后还需要手动分组，因为还是不清晰
	// map分组整理数据
	variantMap := make(map[int]*dto.VariantOptionGroup)

	for _, v := range variantOptions {
		group := variantMap[v.AttributeID]

		if group == nil {
			group = &dto.VariantOptionGroup{
				AttributeID:   v.AttributeID,
				AttributeName: v.AttributeName,
				AttributeCode: v.AttributeCode,
				Options:       []dto.VariantOption{},
			}
			//放进定义的大map里面
			variantMap[v.AttributeID] = group
		}
		//分组options
		//字符串转【】string，dto里是【】
		linkedSkuIDs := strings.Split(v.LinkedSkuIDs, ",")
		//切片直接加数据不用判断，动态变化!
		group.Options = append(group.Options, dto.VariantOption{
			OptionID:     v.OptionID,
			OptionValue:  v.OptionValue,
			OptionCode:   v.OptionCode,
			LinkedSkuIDs: linkedSkuIDs,
		})
	}
	//返回成results，map变【】
	// 打印 variantMap 以供排查

	// 将 map 中的每个 group 添加到 results 中
	for _, group := range variantMap {
		res.Variants = append(res.Variants, *group)
	}

	//存到redis
	saveProductRedis(defaultSkuID, res)

	return res, nil
}

func (P *ProductSkusService) GetUserReviews(ProductCode string, Page int, Limit int, Rating int, sort string) (res dto.ReviewListResponse, err error) {
	db := global.GVA_DB
	var summary dto.ReviewSummary
	// 去除 ProductID 的空格,否则报错
	ProductCode = strings.TrimSpace(ProductCode)

	if ProductCode == "" {
		// return res, errors.New("ProductCode 不能为空")
		return dto.ReviewListResponse{}, ErrProductNotFound
	}
	//思路分析：1.先用code查id
	var productID string
	err = db.Table("products").Select("id").Where("product_code = ?", ProductCode).Scan(&productID).Error
	if err != nil {
		return dto.ReviewListResponse{}, err
	}
	if productID == "" {
		return dto.ReviewListResponse{}, ErrProductNotFound
	}
	//2.拿id找summary
	err = db.Table("review_summaries").
		Select(`
    review_summaries.product_id,
    review_summaries.average_rating,
    review_summaries.review_count,
    review_summaries.rating_1_count,
    review_summaries.rating_2_count,
    review_summaries.rating_3_count,
    review_summaries.rating_4_count,
    review_summaries.rating_5_count
	`).
		Where("review_summaries.product_id = ?", productID).
		Scan(&summary).Error
	if err != nil {
		return dto.ReviewListResponse{}, err
	}
	//3.评论详情,感觉可以和页数设置一起写
	//先查个总评论数,这里必须int64，count返回就是
	//讨论sort，然后再order处设置
	// "排序方式（如  'newest', 'oldest', 'highest_rating', 'lowest_rating', 'most_helpful' 等）"
	var orderBy string
	switch sort {
	case "newest":
		orderBy = "product_reviews.created_at DESC"
	case "oldest":
		orderBy = "product_reviews.created_at ASC"
	case "highest_rating":
		orderBy = "product_reviews.rating DESC"
	case "lowest_rating":
		orderBy = "product_reviews.rating ASC"
	case "most_helpful":
		orderBy = "helpful_count DESC, product_reviews.created_at DESC"
	default:
		orderBy = "product_reviews.created_at DESC"
	}

	// 查询总数
	countQuery := db.Table("product_reviews").
		Where("product_id = ? AND status = 'approved'", productID)
	if Rating >= 1 && Rating <= 5 {
		countQuery = countQuery.Where("rating = ?", Rating)
	}
	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return res, err
	}

	// 查询评论内容
	limit := Limit
	page := Page
	offset := (page - 1) * limit
	var reviews []dto.ReviewInfo
	query := db.Table("product_reviews").
		Select(`
			product_reviews.id AS id,
			product_reviews.nickname AS nickname,
			product_reviews.rating AS rating,
			product_reviews.title AS title,
			product_reviews.comment AS comment,
			product_reviews.created_at AS created_at_formatted,
			IFNULL(GROUP_CONCAT(review_images.image_url), '') AS image_urls,
			(SELECT COUNT(*) FROM user_review_helpful_votes WHERE review_id = product_reviews.id) AS helpful_count
		`).
		Joins("LEFT JOIN review_images ON product_reviews.id = review_images.review_id").
		Where("product_reviews.product_id = ? AND product_reviews.status = 'approved'", productID).
		Order(orderBy).
		Group("product_reviews.id").
		Limit(limit).
		Offset(offset)
	if Rating >= 1 && Rating <= 5 {
		query = query.Where("product_reviews.rating = ?", Rating)
	}

	err = query.Scan(&reviews).Error
	if err != nil {
		return res, err
	}

	// Scan(&reviews).Error
	//reviews是切片，不能直接取元素，要遍历
	for i, _ := range reviews {
		//  CreatedAtFormatted 是string，需要先解析为 time.Time
		t, _ := time.Parse(time.RFC3339, reviews[i].CreatedAtFormatted)
		// 使用 fmt.Sprintf 格式化为 "2023年10月26日" 格式
		reviews[i].CreatedAtFormattedStr = fmt.Sprintf("%d年%d月%d日", t.Year(), int(t.Month()), t.Day())
	}
	//进行一系列判断，page和limit，api层设默认，判断rating和sort
	//rating在后面where处进行判断，sort好像也可以

	// 处理 image_urls，将其转换为 []string
	for i, _ := range reviews {
		if reviews[i].ImageUrls != "" {
			reviews[i].RealImageUrls = strings.Split(reviews[i].ImageUrls, ",")
		}
	}
	//根据拿到的reviews进行筛选,遍历

	//根据输入的显示评论布局
	//int64
	totalPages := int((totalCount + int64(Limit) - 1) / int64(Limit))
	pagination := dto.PaginationInfo{
		CurrentPage: Page,
		Limit:       Limit,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
	}
	//塞数据
	res = dto.ReviewListResponse{
		Summary:    &summary,
		Reviews:    reviews,
		Pagination: pagination,
	}
	return res, nil

}

//var ErrProductNotFound = errors.New("product not found")
//需要统一管理

func (P *ProductSkusService) GetUserQandAs(ProductCode string, Page int, Limit int, sort string) (res dto.QAListResponse, err error) {

	db := global.GVA_DB

	// 去除 ProductID 的空格,否则报错
	ProductCode = strings.TrimSpace(ProductCode)

	if ProductCode == "" {
		// return res, errors.New("ProductCode 不能为空")
		return dto.QAListResponse{}, ErrProductNotFound
	}

	//思路分析：
	// 1.先用code查productid
	//2.然后一起查
	var productID string
	err = db.Table("products").Select("id").Where("product_code = ?", ProductCode).Scan(&productID).Error
	if err != nil {
		return dto.QAListResponse{}, err
	}
	if productID == "" {
		return dto.QAListResponse{}, ErrProductNotFound
	}

	var orderBy string
	switch sort {
	case "newest":
		orderBy = "product_questions.created_at DESC"
	case "oldest":
		orderBy = "product_questions.created_at ASC"
	case "most_helpful":
		orderBy = "helpful_count DESC, product_questions.created_at DESC"
	default:
		orderBy = "product_questions.created_at DESC"
	}
	//查询有回答的question总数
	var totalCount int64
	db.Table("product_questions").
		Joins("JOIN question_answers ON question_answers.question_id = product_questions.id").
		Where("product_questions.product_id = ? AND product_questions.status = 'approved' AND question_answers.status = 'approved'", productID).
		Count(&totalCount)

	// 分页
	limit := Limit
	page := Page
	offset := (page - 1) * limit
	// var qaInfo []dto.QAInfo
	//多个问题和回答要用【】

	//定义临时结构体，接收查到的数据
	var results []struct {
		QuestionID        int64  `json:"id"`
		QuestionText      string `json:"question_text"`
		QuestionCreatedAt string `json:"created_at_formatted"`
		AnswerID          int64  `json:"id"`
		AnswererName      string `json:"answerer_name"`
		AnswerText        string `json:"answer_text"`
		HelpfulCount      int    `json:"helpful_count"`
		AnswerCreatedAt   string `json:"created_at_formatted"`
	}
	//直接在sql把时间格式化
	err = db.Table("product_questions").
		Select(`
        product_questions.id AS question_id,
        product_questions.question_text,
        DATE_FORMAT(product_questions.created_at, '%Y年%m月%d日') AS question_created_at,
        question_answers.id AS answer_id,
        question_answers.answerer_name,
        question_answers.answer_text,
        DATE_FORMAT(question_answers.created_at, '%Y年%m月%d日') AS answer_created_at,
	(SELECT COUNT(*) FROM user_answer_helpful_votes WHERE question_answers.id =user_answer_helpful_votes.answer_id) AS helpful_count
    `).
		//不用leftjoin，直接join就不用考虑没有答案的问题了
		Joins("JOIN question_answers ON product_questions.id = question_answers.question_id").
		// Joins("LEFT JOIN user_answer_helpful_votes ON question_answers.id =user_answer_helpful_votes.answer_id").
		Where("product_questions.product_id = ? AND product_questions.status = 'approved' AND question_answers.status = 'approved'", productID).
		Order(orderBy).
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return dto.QAListResponse{}, err
	}

	//拆分塞到dto
	var qaList []dto.QAInfo
	for _, v := range results {
		qa := dto.QAInfo{
			Question: &dto.QuestionInfo{
				ID:                 v.QuestionID,
				QuestionText:       v.QuestionText,
				CreatedAtFormatted: v.QuestionCreatedAt,
			},
			Answer: &dto.AnswerInfo{
				ID:                 v.AnswerID,
				AnswererName:       v.AnswererName,
				AnswerText:         v.AnswerText,
				HelpfulCount:       v.HelpfulCount,
				CreatedAtFormatted: v.AnswerCreatedAt,
			},
		}
		qaList = append(qaList, qa)
	}

	// //根据输入的显示评论布局
	//int64
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))
	pagination := dto.QaPaginationInfo{
		CurrentPage: Page,
		Limit:       Limit,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
	}
	//塞数据
	res = dto.QAListResponse{
		QAList:     qaList,
		Pagination: pagination,
	}
	return res, nil
}

func (P *ProductSkusService) GetProductImages(SkuId string) (res []dto.SKUImageInfo, err error) {
	db := global.GVA_DB

	// 去除 ProductID 的空格,否则报错
	//SkuId = strings.TrimSpace(SkuId)

	if SkuId == "" {
		return []dto.SKUImageInfo{}, ErrProductNotFound
	}
	var results []dto.SKUImageInfo
	err = db.Table("sku_images").
		Select(`
	id,
	main_image_url,
	thumbnail_url,
	alt_text,
	sort_order
	 `).
		Order("sort_order ASC").
		Where("sku_id = ?", SkuId).Scan(&results).Error
	if err != nil {
		return []dto.SKUImageInfo{}, err
	}
	if SkuId == "" {
		return []dto.SKUImageInfo{}, ErrProductNotFound
	}
	//没有显示[]
	if len(results) == 0 {
		results = []dto.SKUImageInfo{}
	}
	return results, nil
}

func (P *ProductSkusService) GetFavouriteSkuList(UserId uint, Page int, Limit int, sort string) (res dto.FavoriteSKUListResponse, err error) {
	db := global.GVA_DB

	// 获取用户收藏的 SKU ID 列表
	var skuIds []string
	err = db.Table("user_favorite_skus").
		Where("user_id = ?", UserId).
		Pluck("sku_id", &skuIds).Error
	if err != nil {
		return res, err
	}
	if len(skuIds) == 0 {
		return res, nil
	}

	// 获取总数
	var totalCount int64
	err = db.Table("user_favorite_skus").
		Where("user_id = ?", UserId).
		Count(&totalCount).Error
	if err != nil {
		return res, err
	}

	// 排序
	orderBy := "user_favorite_skus.created_at DESC"
	if sort == "oldest" {
		orderBy = "user_favorite_skus.created_at ASC"
	}

	// 分页
	offset := (Page - 1) * Limit

	// 获取 SKU 基本信息
	var favoriteSkuInfo []dto.FavoriteSKUInfo
	err = db.Table("product_skus").
		Select(`
			product_skus.id AS sku_id,
			product_skus.product_id,
			products.name AS product_name,
			products.product_code,
			DATE_FORMAT(user_favorite_skus.created_at, '%Y年%m月%d日 %H:%i:%s') AS added_at_formatted
		`).
		Joins("JOIN user_favorite_skus ON product_skus.id = user_favorite_skus.sku_id").
		Joins("LEFT JOIN products ON products.default_sku_id = product_skus.id").
		Where("user_favorite_skus.user_id = ?", UserId).
		Order(orderBy).
		Limit(Limit).
		Offset(offset).
		Scan(&favoriteSkuInfo).Error
	if err != nil {
		return res, err
	}

	// 获取价格信息
	var prices []struct {
		SkuID    string
		Amount   float64
		Type     string
		TypeName string
	}
	err = db.Table("prices").
		Select("sku_id, price AS amount, price_types.type_code AS type, price_types.name AS type_name").
		Joins("JOIN price_types ON prices.price_type_id = price_types.id").
		Where("sku_id IN ?", skuIds).
		Scan(&prices).Error
	if err != nil {
		return res, err
	}
	priceMap := make(map[string]*dto.PriceInfo)
	for _, p := range prices {
		priceMap[p.SkuID] = &dto.PriceInfo{
			Amount:          p.Amount,
			FormattedAmount: fmt.Sprintf("%s円", humanize.Commaf(p.Amount)),
			Type:            p.Type,
			TypeName:        p.TypeName,
		}
	}
	// 	var primaryImage []struct{
	// 		ID      int
	// 	URL     string
	// 	AltText *string
	// 	}
	// 	err = db.Table("sku_images").
	// 	Select("id, main_image_url, alt_text").
	// 	Joins("JOIN price_types ON prices.price_type_id = price_types.id").
	// 	Where("sku_id IN ?", skuIds).
	// 	Scan(&prices).Error
	// if err != nil {
	// 	return res, err
	// }
	// 获取属性信息
	var attributes []dto.AttributeInfo
	err = db.Table("sku_values").
		Select(`
			sku_values.sku_id,
				attributes.id AS attribute_id,
				attributes.name AS attribute_name,
				attribute_options.id AS option_id,
				attribute_options.value AS option_value,
				attribute_options.option_code AS display_value
			`).
		Joins("JOIN attributes ON sku_values.attribute_id = attributes.id").
		Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
		Where("sku_values.sku_id IN ?", skuIds).
		Scan(&attributes).Error
	if err != nil {
		return res, err
	}
	attrMap := make(map[string][]dto.AttributeInfo)
	for _, attr := range attributes {
		attrMap[attr.SkuID] = append(attrMap[attr.SkuID], attr)
	}

	// 将属性和价格信息合并到 favoriteSkuInfo 中
	for i, sku := range favoriteSkuInfo {
		// 处理时间格式
		// favoriteSkuInfo[i].AddedAtFormatted = sku.AddedAt.Format("2006年01月02日")
		// 填充属性
		favoriteSkuInfo[i].Attributes = attrMap[sku.SkuID]
		// 填充价格
		if p, ok := priceMap[sku.SkuID]; ok {
			favoriteSkuInfo[i].Price = p
		}
	}

	// 分页信息
	totalPages := int((totalCount + int64(Limit) - 1) / int64(Limit))
	pagination := dto.PaginationInfo{
		CurrentPage: Page,
		Limit:       Limit,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
	}

	res = dto.FavoriteSKUListResponse{
		Favorites:  favoriteSkuInfo,
		Pagination: pagination,
	}
	return res, nil
}

func (P *ProductSkusService) AddFavouriteSku(UserId uint, SkuId string) error {
	// 1. 查询 product_skus 表，看 sku_id 是否存在
	//建议用count
	var count int64
	db := global.GVA_DB
	err := db.Table("product_skus").
		Select("product_skus.id").
		Where("id = ?", SkuId).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		// return res, errors.New("ProductCode 不能为空")
		return ErrProductNotFound
	}
	// 2. 查询 user_favorite_skus 表是否已存在记录
	var existing int64
	err = db.Table("user_favorite_skus").
		Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Count(&existing).Error
	if err != nil {
		return err
	}
	if existing > 0 {
		return ErrAlreadyFavorited
	}

	// 3. 插入收藏记录
	err = db.Create(&struct {
		UserId uint   `gorm:"column:user_id"`
		SkuId  string `gorm:"column:sku_id"`
	}{
		UserId: UserId,
		SkuId:  SkuId,
	}).Error

	return err
}

func (P *ProductSkusService) DeleteFavouriteSku(UserId uint, SkuId string) error {
	db := global.GVA_DB

	// 先检查 SKU 是否存在
	var skuCount int64
	if err := db.Table("product_skus").
		Where("id = ?", SkuId).
		Count(&skuCount).Error; err != nil {
		return err
	}
	if skuCount == 0 {
		return ErrProductNotFound
	}

	// 检查是否已收藏
	var favCount int64
	if err := db.Table("user_favorite_skus").
		Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Count(&favCount).Error; err != nil {
		return err
	}
	if favCount == 0 {
		return ErrFavoriteNotFound
	}

	// 执行删除
	if err := db.Table("user_favorite_skus").
		Where("user_id = ? AND sku_id = ?", UserId, SkuId).
		Delete(nil).Error; err != nil {
		return err
	}

	return nil
}
