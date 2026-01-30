package router

import (
	"gin-admin-server/global"

	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	if global.GNA_CONFIG.System.Env == "release" {
		gin.SetMode(gin.ReleaseMode) //DebugMode ReleaseMode TestMode
	} else {
		gin.SetMode(gin.DebugMode)
	}

	Router := gin.New()

	RouterCfg.SetupMiddleware(Router)
	// 为文件提供静态地址
	RouterCfg.SetupStaticFiles(Router)
	//注册路由
	RouterGroup.SetRoutesGroup(Router)

	RouterCfg.SetupErrorHandlers(Router)

	global.GNA_LOG.Info("路由注册成功")
	return Router
}
