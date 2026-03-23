package department

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
)

type _departmentService struct{}

var DepartmentService = new(_departmentService)

// GetDepartmentList 获取全部部门列表（扁平结构，前端自行构建树）
func (ds *_departmentService) GetDepartmentList(c *gin.Context) {
	var list []*SysDepartment
	err := global.GNA_DB.Order("parent_id ASC, sort ASC").Find(&list).Error
	if err != nil {
		global.GNA_LOG.Error("获取部门列表失败: " + err.Error())
		response.FailWithMessage("获取部门列表失败", c)
		return
	}
	response.OkWithData(list, c)
}

// UpdateDepartment 新增或修改部门
func (ds *_departmentService) UpdateDepartment(c *gin.Context) {
	var dept SysDepartment
	err := c.ShouldBindJSON(&dept)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, dept)
		global.GNA_LOG.Error(errMessage)
		response.FailWithMessage(errMessage, c)
		return
	}

	err = global.GNA_DB.Save(&dept).Error
	if err != nil {
		global.GNA_LOG.Error("部门保存失败：" + err.Error())
		response.FailWithMessage("部门保存失败", c)
		return
	}

	response.Ok(c)
}

// DeleteDepartment 删除部门
func (ds *_departmentService) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}

	// 检查是否有子部门
	var count int64
	global.GNA_DB.Model(&SysDepartment{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		response.FailWithMessage("存在子部门，无法删除", c)
		return
	}

	err := global.GNA_DB.Where("id = ?", id).Delete(&SysDepartment{}).Error
	if err != nil {
		global.GNA_LOG.Error("删除部门失败：" + err.Error())
		response.FailWithMessage("删除部门失败", c)
		return
	}

	response.OkWithMessage("删除部门成功", c)
}
