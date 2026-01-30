package router

import (
	"gin-admin-server/api/login"
	"gin-admin-server/global"
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
		r.SetupPublicRoutes(PublicGroup)
	}

	//PrivateGroup := Router.Group(global.GNA_CONFIG.Router.RouterPrefix)
	//PrivateGroup.Use(middleware.JWTAuth())
	{

	}
}

func (r *_routerGroup) SetupPublicRoutes(PublicGroup *gin.RouterGroup) {
	login.LoginRouter.InitLoginRouter(PublicGroup)
}
