package dashboard

import (
	"gin-admin-server/api/department"
	"gin-admin-server/api/dict"
	"gin-admin-server/api/menu"
	"gin-admin-server/api/role"
	"gin-admin-server/api/user"
	"gin-admin-server/global"
	"gin-admin-server/model/response"

	"github.com/gin-gonic/gin"
)

type _dashboardService struct{}

var DashboardService = new(_dashboardService)

func (s *_dashboardService) Stats(c *gin.Context) {
	var userCount, roleCount, menuCount, deptCount, dictTypeCount int64
	_ = global.GNA_DB.Model(&user.SysUser{}).Count(&userCount).Error
	_ = global.GNA_DB.Model(&role.SysRole{}).Count(&roleCount).Error
	_ = global.GNA_DB.Model(&menu.SysMenu{}).Count(&menuCount).Error
	_ = global.GNA_DB.Model(&department.SysDepartment{}).Count(&deptCount).Error
	_ = global.GNA_DB.Model(&dict.SysDictType{}).Count(&dictTypeCount).Error
	response.OkWithData(gin.H{
		"userCount":     userCount,
		"roleCount":     roleCount,
		"menuCount":     menuCount,
		"deptCount":     deptCount,
		"dictTypeCount": dictTypeCount,
	}, c)
}
