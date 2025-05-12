package product

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	ProductRouter
	// CartRouter
}

var (
	getSkuReqApi = api.ApiGroupApp.ProductApiGroup.GetSkuReqApi
	// getUserCartReqApi = api.ApiGroupApp.ProductApiGroup.GetUserCartReqApi
)
