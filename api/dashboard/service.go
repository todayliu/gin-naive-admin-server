package dashboard

import (
	"gin-admin-server/api/department"
	"gin-admin-server/api/dict"
	"gin-admin-server/api/menu"
	"gin-admin-server/api/role"
	"gin-admin-server/api/user"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"

	"github.com/gin-gonic/gin"
)

type _dashboardService struct{}

var DashboardService = new(_dashboardService)

// Stats 仪表盘统计（用户数、角色数等）
// @Summary     仪表盘统计
// @Tags        仪表盘
// @Produce     json
// @Security    AccessToken
// @Success     200 {object} response.Response
// @Router      /dashboard/stats [get]
func (s *_dashboardService) Stats(c *gin.Context) {
	var userCount, roleCount, menuCount, deptCount, dictTypeCount int64
	_ = dbctx.Use(c).Model(&user.SysUser{}).Count(&userCount).Error
	_ = dbctx.Use(c).Model(&role.SysRole{}).Count(&roleCount).Error
	_ = dbctx.Use(c).Model(&menu.SysMenu{}).Count(&menuCount).Error
	_ = dbctx.Use(c).Model(&department.SysDepartment{}).Count(&deptCount).Error
	_ = dbctx.Use(c).Model(&dict.SysDictType{}).Count(&dictTypeCount).Error
	response.OkWithData(gin.H{
		"userCount":     userCount,
		"roleCount":     roleCount,
		"menuCount":     menuCount,
		"deptCount":     deptCount,
		"dictTypeCount": dictTypeCount,
	}, c)
}
