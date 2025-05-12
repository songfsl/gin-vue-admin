package middleware

import (
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)

		// 获取请求的 PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))

		// 判断 Casbin 是否正确加载
		e := casbinService.Casbin()
		if e == nil {
			// 如果 Casbin 对象为空，返回错误
			response.FailWithDetailed(gin.H{}, "Casbin 加载失败，请检查配置", c)
			c.Abort()
			return
		}

		// 判断策略中是否存在权限
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			// 如果权限验证失败，返回权限不足
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}

		// 通过权限检查，继续处理请求
		c.Next()
	}
}
