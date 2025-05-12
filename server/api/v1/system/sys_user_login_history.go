package system

// import (
// 	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
// 	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
// 	"github.com/gin-gonic/gin"
// )

// import ""
type UserLoginoryApi struct{}

// @Tags UserLoginHistory
// @Summary Get login time by user ID
// @Description Get the login time by user ID
// // @Param id query int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /user/getLoginTimeByUsersId [get]
// func (s *UserLoginHistoryApi) GetLoginTimeByUsersId(c *gin.Context) {

// 	var req systemReq.GetLoginTimeByUsersIdReq
// 	err := c.ShouldBindQuery(&req)
// 	if err != nil {
// 		response.FailWithMessage(err.Error(), c)
// 		return
// 	}
// }

// 	list, total, err := operationRecordService.GetSysOperationRecordInfoList(pageInfo)
// 	if err != nil {
// 		global.GVA_LOG.Error("获取失败!", zap.Error(err))
// 		response.FailWithMessage("获取失败", c)
// 		return
// 	}
// 	response.OkWithDetailed(response.PageResult{
// 		List:     list,
// 		Total:    total,
// 		Page:     pageInfo.Page,
// 		PageSize: pageInfo.PageSize,
// 	}, "获取成功", c)
// }
