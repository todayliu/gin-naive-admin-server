package sysconfig

import "github.com/gin-gonic/gin"

type _sysConfigRouter struct{}

var SysConfigRouter = new(_sysConfigRouter)

// InitPublicRouter 注册无需登录的站点展示接口（与私有 config 组路径前缀一致，便于前端统一 /api/config）
func (r *_sysConfigRouter) InitPublicRouter(Public *gin.RouterGroup) {
	Public.GET("config/site-display", SysConfigService.SiteDisplay)
}

func (r *_sysConfigRouter) InitSysConfigRouter(Router *gin.RouterGroup) {
	g := Router.Group("config")
	{
		g.GET("list", SysConfigService.List)
		g.GET("export", SysConfigService.ExportConfigs)
		g.GET("import-template", SysConfigService.DownloadConfigImportTemplate)
		g.POST("import", SysConfigService.ImportConfigs)
		g.POST("add", SysConfigService.Add)
		g.PUT("edit", SysConfigService.Edit)
	}
}
