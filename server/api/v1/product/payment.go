package product

import (
	"errors"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/product"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetPaymentReq struct {
	SkuId     string `json:"skuId"`
	ProductId string `json:"productId"`
}

type GetPaymentReqApi struct{}

var GetPaymentReqApp = new(GetPaymentReqApi)

// GetPaymentMethod 处理获取支付方式的请求
// @Summary 获取用户的支付方式列表
// @Description 根据用户的支付方式获取所有活跃的支付方式，并按排序顺序返回
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.PaymentMethodInfo} "获取支付方式列表成功"
// @Failure 400 {object} response.Response{msg=string} "请求失败"
// @Failure 500 {object} response.Response{msg=string} "内部服务器错误"
// @Router /api/v1/payments/methods [get]
func (g *GetPaymentReqApi) GetPaymentMethod(c *gin.Context) {
	Req, err := product.ProductUserApp.GetPaymentMethod()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		response.FailWithCode("INVALID_PARAMETER", "不正なSKU ID形式です。", c)
		return
	}

	// 返回成功响应
	response.OkWithData(Req, c)
}

// SelectCoupon 商品を追加する
// @Summary 添加商品
// @Description 用户将指定 SKU 添加到自己的カート中，调用此接口必须携带 Authorization: Bearer Token
// @Tags SelectCoupon
// @Accept json
// @Produce json
// @Param sku_id query string true "SKU ID"
// @Param quantity query int true "QUANTITY"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求失败或参数错误"
// @Failure 401 {object} response.Response "未授权，Token 无效或缺失"
// @Router /api/v1/cart/items [post]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) SelectCoupon(c *gin.Context) {
	//1.query获取skuid
	SkuId := c.Query("sku_id")
	if SkuId == "" {
		response.FailWithCode("INVALID_PARAMETER", "sku_idは必須です。", c)
		return
	}
	//2.上下文或许userid
	UserId := utils.GetUserID(c)
	//3.query获取quantity
	//service层对quantity进行了限制，api端不用再重复判断
	quantity := c.Query("quantity")
	var err error
	quantityInt := 1
	if quantity != "" {
		quantityInt, err = strconv.Atoi(quantity)
		if err != nil || quantityInt < 1 {
			response.FailWithCode("INVALID_PARAMETER", "quantityパラメータは1以上の数値で指定してください。", c)
			return
		}
	}
	//还要判断库存情况

	err = product.ProductSkusApp.AddItemsIntoCart(UserId, SkuId, quantityInt)
	if errors.Is(err, product.ErrProductNotFound) {
		// sku_id 不存在或无效
		response.FailWithCode("NOT_FOUND", "SKUが見つかりません。", c)
		return
	}
	if err != nil {
		// sku_id 不存在或无效
		global.GVA_LOG.Error("追加失败!", zap.Error(err))
		return
	}
	response.OkWithMessage("カートを追加しました", c)
}
