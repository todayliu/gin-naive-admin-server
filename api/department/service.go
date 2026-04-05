package department

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
)

type _departmentService struct{}

var DepartmentService = new(_departmentService)

// GetDepartmentList 获取全部部门列表（扁平结构，前端自行构建树）
// @Summary     部门列表
// @Tags        部门
// @Produce     json
// @Security    AccessToken
// @Success     200 {object} response.Response
// @Router      /department/list [get]
func (ds *_departmentService) GetDepartmentList(c *gin.Context) {
	var list []*SysDepartment
	err := dbctx.Use(c).Order("parent_id ASC, sort ASC").Find(&list).Error
	if err != nil {
		global.GNA_LOG.Error("获取部门列表失败: " + err.Error())
		response.FailWithMessage("获取部门列表失败", c)
		return
	}
	global.FillAuditDisplayNames(dbctx.Use(c), &list)
	response.OkWithData(list, c)
}

// UpdateDepartment 新增或修改部门
// @Summary     保存部门
// @Tags        部门
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysDepartment true "部门实体"
// @Success     200 {object} response.Response
// @Router      /department/edit [put]
func (ds *_departmentService) UpdateDepartment(c *gin.Context) {
	var dept SysDepartment
	err := c.ShouldBindJSON(&dept)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, dept)
		global.GNA_LOG.Error(errMessage)
		response.FailWithMessage(errMessage, c)
		return
	}

	err = dbctx.Use(c).Save(&dept).Error
	if err != nil {
		global.GNA_LOG.Error("部门保存失败：" + err.Error())
		response.FailWithMessage("部门保存失败", c)
		return
	}

	response.Ok(c)
}

// DeleteDepartment 删除部门
// @Summary     删除部门
// @Tags        部门
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "部门 ID"
// @Success     200 {object} response.Response
// @Router      /department/delete/{id} [delete]
func (ds *_departmentService) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}

	// 检查是否有子部门
	var count int64
	dbctx.Use(c).Model(&SysDepartment{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		response.FailWithMessage("存在子部门，无法删除", c)
		return
	}

	err := dbctx.Use(c).Where("id = ?", id).Delete(&SysDepartment{}).Error
	if err != nil {
		global.GNA_LOG.Error("删除部门失败：" + err.Error())
		response.FailWithMessage("删除部门失败", c)
		return
	}

	response.OkWithMessage("删除部门成功", c)
}
