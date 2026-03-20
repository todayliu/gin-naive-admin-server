package router

import (
	"gin-admin-server/api/dict"
	"gin-admin-server/api/login"
	"gin-admin-server/api/menu"
	"gin-admin-server/api/role"
	"gin-admin-server/api/user"
	"gin-admin-server/global"
	"gin-admin-server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

var RouterGroup = new(_routerGroup)

type _routerGroup struct{}

func (r *_routerGroup) SetRoutesGroup(Router *gin.Engine) {
	PublicGroup := Router.Group(global.GNA_CONFIG.Router.RouterPrefix)
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "服务状态正常")
		})
	}
	{
		r.SetupPublicRouters(PublicGroup)
	}

	PrivateGroup := Router.Group(global.GNA_CONFIG.Router.RouterPrefix)
	PrivateGroup.Use(middleware.JWTAuth())
	{
		r.SetupPrivateRouters(PrivateGroup)
	}
}

func (r *_routerGroup) SetupPublicRouters(PublicGroup *gin.RouterGroup) {
	login.LoginRouter.InitLoginRouter(PublicGroup)
}

func (r *_routerGroup) SetupPrivateRouters(PrivateGroup *gin.RouterGroup) {
	dict.DictRouter.InitDictRouter(PrivateGroup)
	menu.MenuRouter.InitMenuRouter(PrivateGroup)
	role.RoleRouter.InitRoleRouter(PrivateGroup)
	user.UserRouter.InitUserRouter(PrivateGroup)
}
