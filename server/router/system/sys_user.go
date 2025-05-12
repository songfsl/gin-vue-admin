package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("user")
	productSkusRouter := Router.Group("/api/v1")
	//  cartRouter := Router.Group("/api/v1")
	//productSkusRouter := Router.Group("/api/v1/products")
	// productsRouter := Router.Group("/api/v1/products")
	{
		userRouter.POST("admin_register", baseApi.Register)               // 管理员注册账号
		userRouter.POST("changePassword", baseApi.ChangePassword)         // 用户修改密码
		userRouter.POST("setUserAuthority", baseApi.SetUserAuthority)     // 设置用户权限
		userRouter.DELETE("deleteUser", baseApi.DeleteUser)               // 删除用户
		userRouter.PUT("setUserInfo", baseApi.SetUserInfo)                // 设置用户信息
		userRouter.PUT("setSelfInfo", baseApi.SetSelfInfo)                // 设置自身信息
		userRouter.POST("setUserAuthorities", baseApi.SetUserAuthorities) // 设置用户权限组
		userRouter.POST("resetPassword", baseApi.ResetPassword)           // 设置用户权限组
		userRouter.PUT("setSelfSetting", baseApi.SetSelfSetting)          // 用户界面配置
	}
	{
		userRouterWithoutRecord.POST("getUserList", baseApi.GetUserList) // 分页获取用户列表
		userRouterWithoutRecord.GET("getUserInfo", baseApi.GetUserInfo)  // 获取自身信息
		userRouterWithoutRecord.POST("getLoginHistoryByIdAndTimeRange", baseApi.GetLoginHistoryByIdAndTimeRange)
		productSkusRouter.GET("getTargetProductSkus", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetTargetProductSkus)
		productSkusRouter.GET("getVariantOptions", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetVariantOptions)
		productSkusRouter.GET("getUserReviews", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetUserReviews)
		productSkusRouter.GET("getUserQandAs", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetUserQandAs)
		// productSkusRouter.GET("getProductImages", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetProductImages)

		productSkusRouter.GET("/:SkuId/images", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetProductImages)
		productSkusRouter.GET("/favorites/skus", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetFavouriteSkuList)
		productSkusRouter.POST("/favorites/skus/:sku_id", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.AddFavouriteSku)
		productSkusRouter.DELETE("/favorites/skus/:sku_id", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.DeleteFavouriteSku)
		productSkusRouter.GET("/product/:product_code/related", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetRelatedProductAndCategory)
		productSkusRouter.GET("/product/:product_code/coordinates", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetStaffCoordinate)

		productSkusRouter.POST("/history/viewed-skus/:sku_id", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.AddViewedSkus)
		productSkusRouter.GET("/history/viewed-skus", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetViewedHistory)

		productSkusRouter.POST("cart/items", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.AddItemsIntoCart)
		productSkusRouter.GET("cart", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetCartItems)
		productSkusRouter.DELETE("/cart/items/:sku_id", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.DeleteItemsFromCart)
		productSkusRouter.PUT("/cart/items/:sku_id", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.ChangeItemsInCart)
		// cartRouter.POST("cart/items",v1.ApiGroupApp.
		productSkusRouter.POST("/shipping-addresses", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.CreateShippingAddress)
		productSkusRouter.DELETE("/shipping-addresses/:address_id", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.DeleteShippingAddress)
		productSkusRouter.GET("/shipping-addresses", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.GetShippingAddress)
		productSkusRouter.PUT("/shipping-addresses/:address_id", v1.ApiGroupApp.ProductApiGroup.GetSkuReqApi.ChangeShippingAddress)
	}
}
