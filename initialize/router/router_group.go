package router

import (
	"gin-admin-server/api/dashboard"
	"gin-admin-server/api/department"
	"gin-admin-server/api/dict"
	"gin-admin-server/api/file"
	"gin-admin-server/api/log"
	"gin-admin-server/api/login"
	"gin-admin-server/api/menu"
	"gin-admin-server/api/permissionapi"
	"gin-admin-server/api/position"
	"gin-admin-server/api/profile"
	"gin-admin-server/api/role"
	"gin-admin-server/api/sysconfig"
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
	PrivateGroup.Use(middleware.OperLog())
	{
		r.SetupPrivateRouters(PrivateGroup)
	}
}

func (r *_routerGroup) SetupPublicRouters(PublicGroup *gin.RouterGroup) {
	login.LoginRouter.InitLoginRouter(PublicGroup)
}

func (r *_routerGroup) SetupPrivateRouters(PrivateGroup *gin.RouterGroup) {
	permissionapi.PermissionAPIRouter.InitPermissionAPIRouter(PrivateGroup)
	profile.ProfileRouter.InitProfileRouter(PrivateGroup)
	file.FileRouter.InitFileRouter(PrivateGroup)
	dashboard.DashboardRouter.InitDashboardRouter(PrivateGroup)
	log.LogRouter.InitLogRouter(PrivateGroup)
	sysconfig.SysConfigRouter.InitSysConfigRouter(PrivateGroup)
	department.DepartmentRouter.InitDepartmentRouter(PrivateGroup)
	dict.DictRouter.InitDictRouter(PrivateGroup)
	position.PositionRouter.InitPositionRouter(PrivateGroup)
	menu.MenuRouter.InitMenuRouter(PrivateGroup)
	role.RoleRouter.InitRoleRouter(PrivateGroup)
	user.UserRouter.InitUserRouter(PrivateGroup)
}
