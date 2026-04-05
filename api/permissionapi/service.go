package permissionapi

import (
	"gin-admin-server/api/menu"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/jwt_util"

	"github.com/gin-gonic/gin"
)

type _permissionAPIService struct{}

var PermissionAPIService = new(_permissionAPIService)

// ButtonCodes 返回当前用户在菜单表中有授权的 perms（用户→角色→sys_menu）；超管为库中全部非空 perms。
// @Summary     当前用户按钮权限码
// @Description data.codes 来自 sys_menu.perms，与前端 hasPerm 一致。
// @Tags        权限
// @Produce     json
// @Security    AccessToken
// @Success     200 {object} response.Response
// @Router      /permission/button-codes [get]
func (s *_permissionAPIService) ButtonCodes(c *gin.Context) {
	uid := jwt_util.GetUserID(c)
	if uid == 0 {
		response.FailWithMessage("未登录", c)
		return
	}
	response.OkWithData(gin.H{
		"codes": menu.LoadUserButtonPermCodes(uid),
	}, c)
}
