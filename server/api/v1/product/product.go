package product

import (
	"errors"
	"fmt"

	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/product"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetSkuReq struct {
	SkuId     string `json:"skuId"`
	ProductId string `json:"productId"`
}

type GetSkuReqApi struct{}

var GetSkuReqApiApp = new(GetSkuReqApi)
var (
	ErrProductNotFound  = errors.New("Product not found")
	ErrAlreadyFavorited = errors.New("Product already favorited")
)

// GetTargetProductSkus 获取目标商品的SKU信息
// @Tags 商品SKU
// @Summary 获取目标商品的SKU信息（根据skuId和productId）
// @Description 通过skuId和productId获取相关的SKU信息，返回对应的商品SKU数据。
// @Accept json
// @Produce json
// @Param skuId query string true "SKU ID" // SKU的ID，作为查询参数传递
// @Param productId query string true "商品 ID" // 商品的ID，作为查询参数传递
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "成功返回SKU信息"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /productSkus/getTargetProductSkus [get]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) GetTargetProductSkus(c *gin.Context) {
	skuId := c.Query("skuId")
	productId := c.Query("productId")

	// 校验参数
	if skuId == "" || productId == "" {
		global.GVA_LOG.Error("参数绑定失败!")
		// response.FailWithMessage("skuId和productId不能为空", c)
		response.FailWithCode("INVALID_PARAMETER", "不正な商品識別子です。", c)
		return
	}

	// 获取SKU信息
	Req, err := product.ProductSkusApp.GetTargetProductSkus(skuId, productId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		response.FailWithCode("INVALID_PARAMETER", "不正なSKU ID形式です。", c)
		return
	}

	// 返回成功响应
	response.OkWithDetailed(gin.H{"ProductSkus": Req}, "获取成功", c)
}

// GetVariantOptions 获取变体选项
// @Tags ProductSku
// @Summary 获取指定商品ID对应的变体选项
// @Param productId query string true "商品ID"
// @Success 200 {object} response.Response{data=dto.ProductVariantResponse}
// @Failure 400 {object} response.Response
// @Router /productSkus/getVariantOptions [get]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) GetVariantOptions(c *gin.Context) {

	// 获取查询参数
	productId := c.Query("productId")

	// 检查参数是否为空
	if productId == "" {
		global.GVA_LOG.Error("参数绑定失败!")
		// response.FailWithMessage("productId不能为空", c)
		response.FailWithCode("INVALID_PARAMETER", "不正なSKU ID形式です。", c)
		return
	}
	// 调用服务层方法获取变体选项
	Req, err := product.ProductSkusApp.GetVariantOptions(productId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)

		response.FailWithCode("INVALID_PARAMETER", "不正な商品識別子です。", c)
		return
	}
	fmt.Println("Returning Response Data:", Req)
	// 确保返回的只有一个响应体
	// response.OkWithDetailed(Req, "获取成功", c)
	response.OkWithDetailed(gin.H{"GetVariantOptions": Req}, "获取成功", c)
}

// GetUserReviews
// @Summary 获取指定产品的用户评论信息
// @Description 根据 productCode 获取该商品的评论评分概况（平均评分、各评分数量等）及评论详情列表，支持分页、评分过滤、排序功能。
// @Tags GetUserReviews
// @Accept json
// @Produce json
// @Param productCode query string true "ProductCode"
// @Param page query int false " 取得するページ番号 (1始まり)"
// @Param limit query int false " 1ページあたりのレビュー件数"
// @Param sort query string false "ソート順序('newest', 'oldest', 'highest_rating', 'lowest_rating', 'most_helpful' のいずれかを指定)"
// @Param rating query int false " 指定した評価（星の数、例: 5）"
// @Success 200 {object} response.Response{data=map[string]interface{}} "返回评论列表和评分概况"
// @Failure 400 {object} response.Response "参数错误或查询失败"
// @Router /productSkus/getUserReviews [get]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) GetUserReviews(c *gin.Context) {
	productCode := c.Param("productCode")
	// 获取查询参数

	page := c.Query("page")
	limit := c.Query("limit")
	rating := c.Query("rating")
	sort := c.Query("sort")
	//设置默认值，以后集中管理“写死”的东西
	var err error
	//增加大判断不是空，再转int
	//query默认都是string，所以需要转int
	pageInt := 1
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			response.FailWithCode("INVALID_PARAMETER", "pageパラメータは1以上の数値で指定してください。", c)
			return
		}
	}

	limitInt := 10
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil || limitInt < 1 || limitInt > 100 {
			response.FailWithCode("INVALID_PARAMETER", "limitパラメータは1から100の間で指定してください。", c)
			return
		}
	}

	ratingInt := 0
	if rating != "" {
		ratingInt, err = strconv.Atoi(rating)
		if err != nil || ratingInt < 1 || ratingInt > 5 {
			response.FailWithCode("INVALID_PARAMETER", "ratingパラメータは1から5の間で指定してください。", c)
			return
		}
	}

	// 默认排序为 "newest"
	if sort == "" {
		sort = "newest"
	}
	validSortOptions := []string{"newest", "oldest", "highest_rating", "lowest_rating", "most_helpful"}
	isValidSort := false
	for _, option := range validSortOptions {
		if sort == option {
			isValidSort = true
			break
		}
	}
	if !isValidSort {
		response.FailWithCode("INVALID_PARAMETER", "不正なsortパラメータです。('newest', 'oldest', 'highest_rating', 'lowest_rating', 'most_helpful' のいずれかを指定)", c)
		return
	}
	// 检查 productCode 是否为空
	if productCode == "" {
		response.FailWithCode("INVALID_PARAMETER", "不正な商品識別子です", c)
		return
	}

	// 校验 productCode 长度
	if len(productCode) < 7 {
		response.FailWithCode("INVALID_PARAMETER", "不正な商品識別子です。商品コードは7桁以上で指定してください。", c)
		return
	}

	// 调用服务层方法获取变体选项!!
	Req, err := product.ProductSkusApp.GetUserReviews(productCode, pageInt, limitInt, ratingInt, sort)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		if errors.Is(err, product.ErrProductNotFound) {
			response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
			return
		}
		response.FailWithCode("INVALID_PARAMETER", "不正な商品識別子です。", c)
		return
	}
	fmt.Println("Returning Response Data:", Req)
	// 确保返回的只有一个响应体
	// response.OkWithDetailed(Req, "获取成功", c)
	response.OkWithDetailed(gin.H{"GetUserReviews": Req}, "获取成功", c)
}

// GetUserQandAs
// @Summary 获取指定产品的用户问题以及回答
// @Tags GetUserQandAs
// @Accept json
// @Produce json
// @Param productCode query string false "ProductCode"
// @Param page query int false " 取得するページ番号 (1始まり)"
// @Param limit query int false " 1ページあたりのレビュー件数"
// @Param sort query string false "ソート順序('newest', 'oldest', 'most_helpful' のいずれかを指定)"
// @Success 200 {object} response.Response{data=map[string]interface{}} "返回评论列表和评分概况"
// @Failure 400 {object} response.Response "参数错误或查询失败"
// @Router /productSkus/getUserQandAs [get]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) GetUserQandAs(c *gin.Context) {
	productCode := c.Query("productCode")

	// 获取查询参数

	page := c.Query("page")
	limit := c.Query("limit")
	sort := c.Query("sort")
	//设置默认值，以后集中管理“写死”的东西
	var err error
	//增加大判断不是空，再转int
	//query默认都是string，所以需要转int！！
	pageInt := 1
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			response.FailWithCode("INVALID_PARAMETER", "pageパラメータは1以上の数値で指定してください。", c)
			return
		}
	}

	limitInt := 10
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil || limitInt < 1 || limitInt > 100 {
			response.FailWithCode("INVALID_PARAMETER", "limitパラメータは1から100の間で指定してください。", c)
			return
		}
	}

	// 默认排序为 "newest"
	if sort == "" {
		sort = "newest"
	}
	validSortOptions := []string{"newest", "oldest", "most_helpful"}
	isValidSort := false
	for _, option := range validSortOptions {
		if sort == option {
			isValidSort = true
			break
		}
	}
	if !isValidSort {
		response.FailWithCode("INVALID_PARAMETER", "不正なsortパラメータです。('newest', 'oldest', 'most_helpful' のいずれかを指定)", c)
		return
	}
	// 检查 productCode 是否为空
	if productCode == "" {
		response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		return
	}

	// 校验 productCode 长度
	if len(productCode) < 7 {
		response.FailWithCode("INVALID_PARAMETER", "不正な商品識別子です。商品コードは7桁以上で指定してください。", c)
		return
	}

	// 调用服务层方法获取变体选项!!
	Req, err := product.ProductSkusApp.GetUserQandAs(productCode, pageInt, limitInt, sort)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		if errors.Is(err, product.ErrProductNotFound) {

			response.FailWithCode("INVALID_PARAMETER", "不正な商品識別子です。", c)

			return
		}
		response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		return
	}
	fmt.Println("Returning Response Data:", Req)
	// 确保返回的只有一个响应体
	// response.OkWithDetailed(Req, "获取成功", c)
	response.OkWithDetailed(gin.H{"GetUserQandAs": Req}, "获取成功", c)

}

// GetProductImages 获取指定 SKU 的所有图片信息
// @Tags       GetProductImages
// @Summary    根据 SKU ID 获取商品图片列表
// @Param      SkuId   path      string  true  "SKU ID"
// @Produce    json
// @Success    200  {object}  response.Response{data=map[string]interface{}, msg=string}  "获取成功"
// @Failure    400  {object}  response.Response{msg=string}  "不正な商品識別子です。"
// @Router     /api/v1/products/{SkuId}/images [get]
// @Security   ApiKeyAuth
func (g *GetSkuReqApi) GetProductImages(c *gin.Context) {

	// 获取查询参数
	SkuId := c.Param("SkuId")

	// 检查参数是否为空
	if SkuId == "" {
		global.GVA_LOG.Error("参数绑定失败!")
		// response.FailWithMessage("productId不能为空", c)
		response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		return
	}
	// 调用服务层方法获取变体选项
	Req, err := product.ProductSkusApp.GetProductImages(SkuId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)

		response.FailWithMes("INVALID_PARAMETER", "不正な商品識別子です。", c)
		return
	}
	fmt.Println("Returning Response Data:", Req)
	// 确保返回的只有一个响应体
	response.OkWithDetailed(Req, "获取成功", c)
	// response.OkWithDetailed(gin.H{"GetVariantOptions": Req}, "获取成功", c)
}

// AddFavouriteSku 添加收藏的SKU
// @Summary 添加收藏的SKU
// @Description 用户将指定 SKU 添加到自己的收藏夹中，调用此接口必须携带 Authorization: Bearer Token
// @Tags AddFavouriteSku
// @Accept json
// @Produce json
// @Param sku_id path string true "SKU ID"
// @Success 200 {object} response.Response "收藏成功"
// @Failure 400 {object} response.Response "请求失败或参数错误"
// @Failure 401 {object} response.Response "未授权，Token 无效或缺失"
// @Router /api/v1/favorites/skus/{sku_id} [post]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) AddFavouriteSku(c *gin.Context) {
	//路径获取
	SkuId := c.Param("sku_id")
	if SkuId == "" {
		response.FailWithCode("INVALID_PARAMETER", "sku_idは必須です。", c)
		return
	}

	UserId := utils.GetUserID(c)

	err := product.ProductSkusApp.AddFavouriteSku(UserId, SkuId)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			// sku_id 不存在或无效
			response.FailWithCode("NOT_FOUND", "SKUが見つかりません。", c)
			return
		}
		if errors.Is(err, product.ErrAlreadyFavorited) {
			// 已收藏
			response.FailWithCode("ALREADY_EXISTS", "既にお気に入りに追加済みです。", c)
			return
		}
	}
	response.OkWithMessage("お気に入りに追加しました", c)
}

// DeleteFavouriteSku 删除用户收藏的 SKU
// @Summary 删除用户收藏的 SKU
// @Description 根据用户的 SKU ID 删除其收藏的 SKU
// @Tags DeleteFavouriteSku
// @Accept json
// @Produce json
// @Param sku_id path string true "SKU ID" // 路径参数，SKU的唯一标识
// @Success 200 {object} response.Response "删除成功的响应"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "删除失败的响应"
// @Router /api/v1/favorites/skus/{sku_id} [delete]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) DeleteFavouriteSku(c *gin.Context) {
	//路径获取
	SkuId := c.Param("sku_id")

	UserId := utils.GetUserID(c) // 从JWT中获取userID

	err := product.ProductSkusApp.DeleteFavouriteSku(UserId, SkuId)
	if err != nil {
		switch err {
		case product.ErrProductNotFound:
			response.FailWithCode("SKU_NOT_FOUND", "指定された商品が存在しません。", c)
		case product.ErrFavoriteNotFound:
			response.FailWithCode("NOT_FOUND", "お気に入りが見つかりません。", c)
		default:
			response.FailWithMessage("削除に失敗しました", c)
		}
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetFavouriteSkuList 获取当前认证用户的收藏 SKU 列表（带分页和排序）
// @Summary 获取用户收藏SKU列表
// @Description 获取当前用户收藏的SKU信息，支持分页与排序
// @Tags GetFavouriteSkuList
// @Accept json
// @Produce json
// @Param page query int false " 取得するページ番号 (1始まり)"
// @Param limit query int false " 1ページあたりのレビュー件数"
// @Param sort query string false "ソート順序('newest', 'oldest', のいずれかを指定)"
// @Success 200 {object} response.Response
// @Router /api/v1/favorites/skus [get]
// @Security   ApiKeyAuth
func (g *GetSkuReqApi) GetFavouriteSkuList(c *gin.Context) {

	//调utils方法获取userid，
	userId := utils.GetUserID(c) // 从JWT中获取userID

	// 获取查询参数
	page := c.Query("page")
	limit := c.Query("limit")
	sort := c.Query("sort")
	//设置默认值，以后集中管理“写死”的东西
	var err error
	//增加大判断不是空，再转int
	//query默认都是string，所以需要转int！！
	pageInt := 1
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			response.FailWithCode("INVALID_PARAMETER", "pageパラメータは1以上の数値で指定してください。", c)
			return
		}
	}

	limitInt := 10
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil || limitInt < 1 || limitInt > 100 {
			response.FailWithCode("INVALID_PARAMETER", "limitパラメータは1から100の間で指定してください。", c)
			return
		}
	}
	// 默认排序为 "newest"
	if sort == "" {
		sort = "newest"
	}
	validSortOptions := []string{"newest", "oldest"}
	isValidSort := false
	for _, option := range validSortOptions {
		if sort == option {
			isValidSort = true
			break
		}
	}
	if !isValidSort {
		response.FailWithMes("INVALID_PARAMETER", "不正なsortパラメータです。('newest', 'oldest', のいずれかを指定)", c)
		return
	}

	// 调用服务层方法获取变体选项!!
	Req, err := product.ProductSkusApp.GetFavouriteSkuList(userId, pageInt, limitInt, sort)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		if errors.Is(err, product.ErrProductNotFound) {

			response.FailWithMes("INVALID_PARAMETER", "不正な商品識別子です。", c)

			return
		}
		// response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		// return
	}

	// 确保返回的只有一个响应体
	response.OkWithDetailed(Req, "获取成功", c)
	//	response.OkWithDetailed(gin.H{"GetUserQandAs": Req}, "获取成功", c)

}

// GetRelatedProductAndCategory
// @Summary     関連商品を取得
// @Description 指定された商品コードに基づいて関連商品情報を取得します（最大10件まで）
// @Tags        GetRelatedProductAndCategory
// @Param       product_code path string true "商品コード（例：ABC123）"
// @Param       limit query int false "取得件数（1〜10）" minimum(1) maximum(10) default(5)
// @Success     200 {object} response.Response  "获取成功"
// @Failure     400 {object} response.Response "パラメータ不正、または商品コードが存在しない"
// @Router      /api/v1/product/{product_code}/related [get]
// @Security   ApiKeyAuth
func (g *GetSkuReqApi) GetRelatedProductAndCategory(c *gin.Context) {
	//路径获取
	ProductCode := c.Param("product_code")
	limit := c.Query("limit")
	var err error
	limitInt := 5
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil || limitInt < 1 || limitInt > 10 {
			response.FailWithMes("INVALID_PARAMETER", "limitパラメータは1から10の間で指定してください。", c)
			return
		}
	}
	// 检查参数是否为空
	if ProductCode == "" {
		global.GVA_LOG.Error("参数绑定失败!")
		// response.FailWithMessage("productId不能为空", c)
		response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		return
	}
	// 调用服务层方法获取变体选项
	Req, err := product.ProductSkusApp.GetRelatedProduct(ProductCode, limitInt)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)

		response.FailWithMes("INVALID_PARAMETER", "不正な商品識別子です。", c)
		return
	}
	// 确保返回的只有一个响应体
	response.OkWithDetailed(Req, "获取成功", c)
	// response.OkWithDetailed(gin.H{"GetVariantOptions": Req}, "获取成功", c)
}

// GetStaffCoordinate
// @Summary     「コーディネートセットの概要」を取得
// @Description 現在閲覧している商品詳細ページに関連する「コーディネートセットの概要」リストを返す
// @Tags     GetStaffCoordinate
// @Param       product_code path string true "商品コード（例：ABC123）"
// @Param       limit query int false "取得件数（1〜10）" minimum(1) maximum(5) default(4)
// @Success     200 {object} response.Response  "获取成功"
// @Failure     400 {object} response.Response "パラメータ不正、または商品コードが存在しない"
// @Router      /api/v1/product/{product_code}/coordinates [get]
// @Security   ApiKeyAuth
func (g *GetSkuReqApi) GetStaffCoordinate(c *gin.Context) {
	//路径获取
	ProductCode := c.Param("product_code")
	limit := c.Query("limit")
	var err error
	limitInt := 4
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil || limitInt < 1 || limitInt > 5 {
			response.FailWithMes("INVALID_PARAMETER", "limitパラメータは1から10の間で指定してください。", c)
			return
		}
	}
	// 检查参数是否为空
	if ProductCode == "" {
		global.GVA_LOG.Error("参数绑定失败!")
		// response.FailWithMessage("productId不能为空", c)
		response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		return
	}
	// 调用服务层方法获取变体选项
	Req, err := product.ProductSkusApp.GetStaffCoordinate(ProductCode, limitInt)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)

		response.FailWithMes("INVALID_PARAMETER", "不正な商品識別子です。", c)
		return
	}
	// 确保返回的只有一个响应体
	response.OkWithDetailed(Req, "获取成功", c)
	// response.OkWithDetailed(gin.H{"GetVariantOptions": Req}, "获取成功", c)
}

// AddViewedSkus 商品閲覧履歴
// @Summary 添加商品閲覧履歴
// @Description 用户将指定 SKU 添加到自己的閲覧履歴中，调用此接口必须携带 Authorization: Bearer Token
// @Tags AddViewedSkus
// @Accept json
// @Produce json
// @Param sku_id path string true "SKU ID"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求失败或参数错误"
// @Failure 401 {object} response.Response "未授权，Token 无效或缺失"
// @Router /api/v1/history/viewed-skus/{sku_id} [post]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) AddViewedSkus(c *gin.Context) {
	//路径获取skuid
	SkuId := c.Param("sku_id")
	if SkuId == "" {
		response.FailWithCode("INVALID_PARAMETER", "sku_idは必須です。", c)
		return
	}
	//上下文或许userid
	UserId := utils.GetUserID(c)

	err := product.ProductSkusApp.AddViewedSkus(UserId, SkuId)
	if errors.Is(err, product.ErrProductNotFound) {
		// sku_id 不存在或无效
		response.FailWithCode("NOT_FOUND", "SKUが見つかりません。", c)
		return
	}
	if errors.Is(err, product.ErrAlreadyFavorited) {
		// 已收藏
		response.FailWithCode("ALREADY_EXISTS", "既に閲覧履歴を記録しました", c)
		return
	}
	response.OkWithMessage("閲覧履歴を記録しました", c)
}

// GetViewedHistory 获取当前认证用户的浏览 SKU 列表（带分页和排序）
// @Summary 获取用户浏览SKU列表
// @Description 获取当前用户浏览的SKU信息，支持分页与排序
// @Tags GetViewedHistory
// @Accept json
// @Produce json
// @Param page query int false " 取得するページ番号 (1始まり)"
// @Param limit query int false " 1ページあたりのレビュー件数"
// @Success 200 {object} response.Response
// @Router /api/v1/history/viewed-skus/ [get]
// @Security   ApiKeyAuth
func (g *GetSkuReqApi) GetViewedHistory(c *gin.Context) {

	//调utils方法获取userid，
	userId := utils.GetUserID(c) // 从JWT中获取userID

	// 获取查询参数
	page := c.Query("page")
	limit := c.Query("limit")
	//设置默认值，以后集中管理“写死”的东西
	var err error
	//增加大判断不是空，再转int
	//query默认都是string，所以需要转int！！
	pageInt := 1
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			response.FailWithCode("INVALID_PARAMETER", "pageパラメータは1以上の数値で指定してください。", c)
			return
		}
	}

	limitInt := 10
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil || limitInt < 1 || limitInt > 100 {
			response.FailWithCode("INVALID_PARAMETER", "limitパラメータは1から100の間で指定してください。", c)
			return
		}
	}
	// 调用服务层方法获取变体选项!!
	Req, err := product.ProductSkusApp.GetViewedHistory(userId, pageInt, limitInt)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		if errors.Is(err, product.ErrProductNotFound) {

			response.FailWithMes("INVALID_PARAMETER", "不正な商品識別子です。", c)

			return
		}
		// response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		// return
	}

	// 确保返回的只有一个响应体
	response.OkWithDetailed(Req, "获取成功", c)
	//	response.OkWithDetailed(gin.H{"GetUserQandAs": Req}, "获取成功", c)

}

// AddItemsIntoCart 商品を追加する
// @Summary 添加商品
// @Description 用户将指定 SKU 添加到自己的カート中，调用此接口必须携带 Authorization: Bearer Token
// @Tags AddItemsIntoCart
// @Accept json
// @Produce json
// @Param sku_id query string true "SKU ID"
// @Param quantity query int true "QUANTITY"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求失败或参数错误"
// @Failure 401 {object} response.Response "未授权，Token 无效或缺失"
// @Router /api/v1/cart/items [post]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) AddItemsIntoCart(c *gin.Context) {
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

// GetCartItems 获取当前认证用户カート的 SKU 列表
// @Summary カート内容取得
// @Description 获取当前用户的カートSKU信息
// @Tags GetCartItems
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/cart [get]
// @Security   ApiKeyAuth
func (g *GetSkuReqApi) GetCartItems(c *gin.Context) {

	//调utils方法获取userid，
	userId := utils.GetUserID(c) // 从JWT中获取userID

	// 调用服务层方法获取变体选项!!
	Req, err := product.ProductSkusApp.GetCartItems(userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		if errors.Is(err, product.ErrProductNotFound) {

			response.FailWithMes("INVALID_PARAMETER", "不正な商品識別子です。", c)

			return
		}
		// response.FailWithCode("NOT_FOUND", "商品が見つかりません。", c)
		// return
	}

	// 确保返回的只有一个响应体
	response.OkWithDetailed(Req, "获取成功", c)
	//	response.OkWithDetailed(gin.H{"GetUserQandAs": Req}, "获取成功", c)

}

// DeleteItemsFromCart カート商品削除
// @Summary 删除用户cart的 SKU
// @Description 根据用户的 SKU ID 删除其cart的 SKU
// @Tags DeleteItemsFromCart
// @Accept json
// @Produce json
// @Param sku_id path string true "SKU ID" // 路径参数，SKU的唯一标识
// @Success 200 {object} response.Response "删除成功的响应"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "删除失败的响应"
// @Router /api/v1/cart/items/{sku_id} [delete]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) DeleteItemsFromCart(c *gin.Context) {
	//路径获取
	SkuId := c.Param("sku_id")

	UserId := utils.GetUserID(c) // 从JWT中获取userID

	err := product.ProductSkusApp.DeleteItemsFromCart(UserId, SkuId)
	if err != nil {
		switch err {
		case product.ErrProductNotFound:
			response.FailWithCode("SKU_NOT_FOUND", "指定された商品が存在しません。", c)
		case product.ErrFavoriteNotFound:
			response.FailWithCode("NOT_FOUND", "お気に入りが見つかりません。", c)
		default:
			response.FailWithMessage("削除に失敗しました", c)
		}
		return
	}

	response.OkWithMessage("削除しました", c)
}

// ChangeItemsInCart 商品数量変更
// @Summary カート商品数量変更
// @Description  sku_id の商品の数量を、リクエストボディで指定された新しい数量に変更する，调用此接口必须携带 Authorization: Bearer Token
// @Tags ChangeItemsInCart
// @Accept json
// @Produce json
// @Param sku_id query string true "SKU ID"
// @Param quantity query int true "QUANTITY"
// @Success 200 {object} response.Response "变更成功"
// @Failure 400 {object} response.Response "请求失败或参数错误"
// @Failure 401 {object} response.Response "未授权，Token 无效或缺失"
// @Router /api/v1/cart/items/{sku_id} [put]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) ChangeItemsInCart(c *gin.Context) {
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
	//sevice端已经判断
	err = product.ProductSkusApp.AddItemsIntoCart(UserId, SkuId, quantityInt)
	if errors.Is(err, product.ErrProductNotFound) {
		// sku_id 不存在或无效
		response.FailWithCode("NOT_FOUND", "SKUが見つかりません。", c)
		return
	}
	if err != nil {
		// sku_id 不存在或无效
		global.GVA_LOG.Error("変更失败!", zap.Error(err))
		return
	}
	response.OkWithMessage("カートを変更しました", c)
}

// CreateShippingAddress 配送先住所を追加する
// @Summary 配送先住所追加
// @Description 用户将指定住所添加到自己的地址簿中，调用此接口必须携带 Authorization: Bearer Token
// @Tags CreateShippingAddress
// @Accept json
// @Produce json
// @Param data body dto.ShippingAddressInput true "配送先住所信息"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求失败或参数错误"
// @Failure 401 {object} response.Response "未授权，Token 无效或缺失"
// @Router /api/v1/shipping-addresses [post]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) CreateShippingAddress(c *gin.Context) {

	//上下文或许userid
	UserId := utils.GetUserID(c)
	var req dto.ShippingAddressInput

	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果绑定失败，返回错误信息
		global.GVA_LOG.Error("绑定失败", zap.Error(err))
		response.FailWithMes("INVALID_PARAMETER", "请求体数据格式错误", c)
		return
	}
	err := product.ProductSkusApp.CreateShippingAddress(UserId, req)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithMes("INVALID_PARAMETER", "不正な配送先住所です。", c)
		return
	}

	response.OkWithMessage("配送先住所を追加しました", c)
}

// DeleteShippingAddress 配送先住所削除
// @Summary 配送先住所削除
// @Description 根据用户的addressid 删除其指定地址
// @Tags DeleteShippingAddress
// @Accept json
// @Produce json
// @Param address_id path int true "ADDRESS ID"
// @Success 200 {object} response.Response "删除成功的响应"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "删除失败的响应"
// @Router /api/v1/shipping-addresses/{address_id} [delete]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) DeleteShippingAddress(c *gin.Context) {

	//路径获取
	AddressId := c.Param("address_id")
	AddressIdInt, _ := strconv.Atoi(AddressId)
	UserId := utils.GetUserID(c) // 从JWT中获取userID

	err := product.ProductSkusApp.DeleteShippingAddress(UserId, AddressIdInt)
	if err != nil {
		switch err {
		case product.ErrProductNotFound:
			response.FailWithCode("SKU_NOT_FOUND", "指定された送先住所が存在しません。", c)
		case product.ErrFavoriteNotFound:
			response.FailWithCode("NOT_FOUND", "お送先住所が見つかりません。", c)
		default:
			response.FailWithMessage("削除に失敗しました", c)
		}
		return
	}

	response.OkWithMessage("削除しました", c)
}

// GetShippingAddress 配送先住所リスト取得
// @Summary 住所リスト取得
// @Description 認証済みユーザーの登録済み全配送先住所リストを取得する。调用此接口必须携带 Authorization: Bearer Token
// @Tags GetShippingAddress
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/shipping-addresses [get]
// @Security   ApiKeyAuth
func (g *GetSkuReqApi) GetShippingAddress(c *gin.Context) {

	//调utils方法获取userid，
	userId := utils.GetUserID(c) // 从JWT中获取userID

	// 调用服务层方法获取变体选项!!
	Req, err := product.ProductSkusApp.GetShippingAddress(userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		// response.FailWithMessage("获取失败", c)
		if errors.Is(err, product.ErrProductNotFound) {

			response.FailWithMes("INVALID_PARAMETER", "不正な住所リストです。", c)

			return
		}
		response.FailWithCode("NOT_FOUND", "住所リストが見つかりません。", c)
		return
	}

	// 确保返回的只有一个响应体
	response.OkWithDetailed(Req, "取得しました", c)
	//	response.OkWithDetailed(gin.H{"GetUserQandAs": Req}, "获取成功", c)

}

// ChangeShippingAddress 配送先住所編集
// @Summary 配送先住所編集
// @Description 用户将指定住所編集，调用此接口必须携带 Authorization: Bearer Token
// @Tags ChangeShippingAddress
// @Accept json
// @Produce json
// @Param address_id path int true "ADDRESS ID"
// @Param data body dto.ShippingAddressInput true "更改配送先住所信息"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求失败或参数错误"
// @Failure 401 {object} response.Response "未授权，Token 无效或缺失"
// @Router /api/v1/shipping-addresses/{address_id} [put]
// @Security ApiKeyAuth
func (g *GetSkuReqApi) ChangeShippingAddress(c *gin.Context) {

	//上下文或许userid
	UserId := utils.GetUserID(c)
	AddressId := c.Param("address_id")
	AddressIdInt, _ := strconv.Atoi(AddressId)
	var req dto.ShippingAddressInput

	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果绑定失败，返回错误信息
		global.GVA_LOG.Error("変更失败", zap.Error(err))
		response.FailWithMes("INVALID_PARAMETER", "请求体数据格式错误", c)
		return
	}
	err := product.ProductSkusApp.ChangeShippingAddress(UserId, AddressIdInt, req)
	if err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMes("INVALID_PARAMETER", "不正な配送先住所です。", c)
		return
	}

	response.OkWithMessage("配送先住所を変更しました", c)
}
