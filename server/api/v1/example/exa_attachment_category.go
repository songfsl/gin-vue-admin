package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	common "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetById struct {
	ID int `json:"id" form:"id"`
}
type AttachmentCategoryApi struct{}

// @Tags      AttachmentCategory
// @Summary   获取媒体库分类列表
// @Security  AttachmentCategory
// @Produce   application/json
// @Success   200 {object} response.Response "操作成功"
// @Failure   400 {object} response.Response "参数错误"
// @Router    /attachmentCategory/getCategoryList [get]
func (a *AttachmentCategoryApi) GetCategoryList(c *gin.Context) {
	res, err := attachmentCategoryService.GetCategoryList()
	if err != nil {
		global.GVA_LOG.Error("获取分类列表失败!", zap.Error(err))
		response.FailWithMessage("获取分类列表失败", c)
		return
	}
	response.OkWithData(res, c)
}

// @Tags      AttachmentCategory
// @Summary   添加媒体库分类
// @Security  AttachmentCategory
// @Accept    application/json
// @Produce   application/json
// @Param     data body example.ExaAttachmentCategory true "媒体库分类数据"
// @Success   200 {object} response.Response{msg=string} "添加成功"
// @Failure   400 {object} response.Response{msg=string} "添加失败"
// @Router    /attachmentCategory/addCategory [post]
func (a *AttachmentCategoryApi) AddCategory(c *gin.Context) {
	var req example.ExaAttachmentCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := attachmentCategoryService.AddCategory(&req); err != nil {
		global.GVA_LOG.Error("创建/更新失败!", zap.Error(err))
		response.FailWithMessage("创建/更新失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("创建/更新成功", c)
}

// @Tags      AttachmentCategory
// @Summary   删除媒体库分类
// @Security  AttachmentCategory
// @Accept    application/json
// @Produce   application/json

// @Param data body common.GetById true "通过ID获取"

// @Success   200 {object} response.Response{msg=string} "删除成功"
// @Failure   400 {object} response.Response{msg=string} "删除失败"
// @Router    /attachmentCategory/deleteCategory [post]
func (a *AttachmentCategoryApi) DeleteCategory(c *gin.Context) {
	var req common.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if req.ID == 0 {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := attachmentCategoryService.DeleteCategory(&req.ID); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
