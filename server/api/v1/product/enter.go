package product

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	GetSkuReqApi
	GetPaymentReqApi
}

var (
	productSkusService = service.ServiceGroupApp.ProductServiceGroup.ProductSkusService
	productUserService = service.ServiceGroupApp.ProductServiceGroup.ProductUserService
)
