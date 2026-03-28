package router

import (
	_ "gin-admin-server/docs" // swag 生成的 Swagger 文档
	"gin-admin-server/global"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute() *gin.Engine {
	if global.GNA_CONFIG.System.Env == "release" {
		gin.SetMode(gin.ReleaseMode) //DebugMode ReleaseMode TestMode
	} else {
		gin.SetMode(gin.DebugMode)
	}

	Router := gin.New()

	RouterCfg.SetupMiddleware(Router)

	// Swagger UI：http://<host>:<port>/swagger/index.html
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	RouterCfg.SetupErrorHandlers(Router)
	// 为文件提供静态地址
	RouterCfg.SetupStaticFiles(Router)
	//注册路由
	RouterGroup.SetRoutesGroup(Router)

	global.GNA_LOG.Info("路由注册成功")
	return Router
}
