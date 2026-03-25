package sysconfig

import "github.com/gin-gonic/gin"

type _sysConfigRouter struct{}

var SysConfigRouter = new(_sysConfigRouter)

func (r *_sysConfigRouter) InitSysConfigRouter(Router *gin.RouterGroup) {
	g := Router.Group("config")
	{
		g.GET("list", SysConfigService.List)
		g.PUT("edit", SysConfigService.Edit)
	}
}
