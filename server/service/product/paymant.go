package product

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ProductUserService struct{}

var ProductUserApp = new(ProductUserService)
var (
	ErrConponNotFound          = errors.New("conpon not found")
	ErrCouponInvalid           = errors.New("conpon not vaild")
	ErrCouponExpired           = errors.New("conpon out of date")
	ErrCouponMinPurchaseNotMet = errors.New("conpon not minimum amount")
	// ErrAlreadyFavorited = errors.New("product already favorited")
	// ErrFavoriteNotFound = errors.New("favorite not found")
)

func (P *ProductUserService) GetPaymentMethod() (req []dto.PaymentMethodInfo, err error) {
	db := global.GVA_DB
	var results []dto.PaymentMethodInfo
	err = db.Table("payment_methods").
		Select(`
		id AS method_id,
		method_code,
		name,
		description,
		is_active,
		sort_order
			`).
		Where("is_active = ?", true).
		Order("sort_order ASC").
		Scan(&results).Error
	if err != nil {
		return []dto.PaymentMethodInfo{}, err
	}
	//没有显示[]
	if len(results) == 0 {
		results = []dto.PaymentMethodInfo{}
	}
	return results, nil

}

func (P *ProductUserService) SelectCoupon(UserId uint, ConponCode string) (res dto.CheckoutInfoResponse, err error) {
	//查询有没有购物车付款记录
	db := global.GVA_DB
	type checkout struct {
		UserId               uint
		UsedPoints           int
		PointsDiscountAmount float64
		ShippingFee          float64
	}
	var session checkout
	err = db.Table("checkout_sessions").Where("user_id = ?", UserId).Scan(&session).Error
	if err != nil {
		return res, err
	}
	// 查询conpons表，看conponcode是否存在,并且active
	//校验conpon,需要自定义结构体，校验多个字段
	var coupon struct {
		CouponID           uint64  `json:"coupon_id"`
		CouponCode         string  `json:"coupon_code"`
		Name               string  `json:"name"`
		Description        *string `json:"description,omitempty"`
		DiscountText       string  `json:"discount_text"`
		DiscountValue      float64 `json:"discount_value"`
		IsActive           bool
		MinPurchaseAccount float64
		MaxPurchaseAccount float64
		StartDate          time.Time
		EndDate            time.Time
	}
	err = db.Table("conpons").
		Select("conpon_code").
		Where("conpon_code = ? ", ConponCode).
		Scan(&coupon).Error
	//开始校验
	if err != nil {
		return res, ErrConponNotFound
	}
	//active
	if !coupon.IsActive {
		return res, ErrCouponInvalid
	}
	if time.Now().After(coupon.EndDate) || time.Now().Before(coupon.StartDate) {
		return res, ErrCouponExpired
	}
	//金额
	if session.PointsDiscountAmount < coupon.MinPurchaseAccount {
		return res, ErrCouponMinPurchaseNotMet
	}
	// var count int64
	// err = db.Table("conpons").
	// 	Select("conpon_code").
	// 	Where("conpon_code = ? AND is _active", ConponCode, true).
	// 	Count(&count).Error
	// if count == 0 {
	// 	return res, ErrConponNotFound
	// }
	// if
	// 2. 查询user表是否有且可用
	var existing int64
	err = db.Table("checkout_sessions").
		Where("user_id = ?", UserId).
		Count(&existing).Error

	// newCheck := checkout{
	// 	UserId:               UserId,
	// 	UsedPoints:           UsedPoints,
	// 	PointsDiscountAmount: PointsDiscountAmount,
	// 	ShippingFee:          ShippingFee,
	// }
	// if existing > 0 {
	// 	err = db.Table("checkout_sessions").
	// 		Where("user_id = ? ", UserId).
	// 		Updates(newCheck).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// // 3. 没有就插入，指定表名！(匿名结构体插入数据gorm无法判断表名)
	// err = db.Table("user_conpons").Create(&struct {
	// 	UserId     uint
	// 	ConponCode string
	// }{
	// 	UserId:     UserId,
	// 	ConponCode: ConponCode,
	// }).Error

	return res, err

}

// func (P *ProductUserService) CheckCouponAndPoint(UserId uint) (res dto.CheckoutInfoResponse, err error) {

// }
