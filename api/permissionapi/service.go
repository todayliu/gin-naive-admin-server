package permissionapi

import (
	"gin-admin-server/model/response"
	"gin-admin-server/permission"
	"gin-admin-server/utils/jwt_util"

	"github.com/gin-gonic/gin"
)

type _permissionAPIService struct{}

var PermissionAPIService = new(_permissionAPIService)

// ButtonCodes 返回当前用户拥有的按钮权限标识（与菜单按钮 perms 一致，超管为全部注册码）
// @Summary     当前用户按钮权限码
// @Description data.codes 为权限字符串数组，与前端 v-permission 一致。
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
	codes := permission.LoadUserEffectivePermCodes(uid)
	response.OkWithData(gin.H{
		"codes": codes,
	}, c)
}
