package product

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/product"
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (s *ProductRouter) InitSkuRouter(Router *gin.RouterGroup) {
	ProductRouter := Router.Group("sku")

	{
		ProductRouter.GET("get", product.GetSkuReqApiApp.GetTargetProductSkus)
		ProductRouter.GET("options", product.GetSkuReqApiApp.GetVariantOptions)
		ProductRouter.GET("images", product.GetSkuReqApiApp.GetProductImages)
		ProductRouter.GET("favorites", product.GetSkuReqApiApp.GetFavouriteSkuList)
		ProductRouter.POST("favorites", product.GetSkuReqApiApp.AddFavouriteSku)
		ProductRouter.DELETE("favorites", product.GetSkuReqApiApp.DeleteFavouriteSku)
		ProductRouter.GET("related", product.GetSkuReqApiApp.GetRelatedProductAndCategory)
		ProductRouter.GET("coordinates", product.GetSkuReqApiApp.GetStaffCoordinate)
		ProductRouter.POST("viewedhistory", product.GetSkuReqApiApp.AddViewedSkus)
		ProductRouter.GET("viewedhistory", product.GetSkuReqApiApp.GetViewedHistory)

		ProductRouter.POST("items", product.GetSkuReqApiApp.AddItemsIntoCart)
		ProductRouter.GET("items", product.GetSkuReqApiApp.GetCartItems)
		ProductRouter.DELETE("items", product.GetSkuReqApiApp.DeleteItemsFromCart)
		ProductRouter.PUT("items", product.GetSkuReqApiApp.ChangeItemsInCart)

		ProductRouter.POST("adresses", product.GetSkuReqApiApp.CreateShippingAddress)
		ProductRouter.DELETE("adresses", product.GetSkuReqApiApp.DeleteShippingAddress)
		ProductRouter.GET("adresses", product.GetSkuReqApiApp.GetShippingAddress)
		ProductRouter.PUT("adresses", product.GetSkuReqApiApp.ChangeShippingAddress)
	}
}
