package router

import (
	"gin-admin-server/global"
	"gin-admin-server/middleware"
	"gin-admin-server/model/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

var RouterCfg = new(_router)

type _router struct{}

func (r *_router) SetupMiddleware(Router *gin.Engine) {
	if global.GNA_CONFIG.System.Env == "release" {
		Router.Use(gin.Recovery())               // 生产环境必须有的恢复中间件
		Router.Use(middleware.SecurityHeaders()) // 安全头
		Router.Use(middleware.CORS())            // 跨域
		Router.Use(middleware.Logger())          //请求日志
	} else {
		//Router.Use(gin.Logger())        // 日志
		Router.Use(gin.Recovery())      // 恢复
		Router.Use(middleware.CORS())   // 跨域
		Router.Use(middleware.Logger()) //请求日志
	}
}

func (r *_router) SetupStaticFiles(Router *gin.Engine) {
	// 静态文件目录
	if global.GNA_CONFIG.Router.Path != "" {
		Router.StaticFS(global.GNA_CONFIG.Router.Path, http.Dir(global.GNA_CONFIG.Router.Path))
	}
}

func (r *_router) SetupErrorHandlers(Router *gin.Engine) {
	Router.NoRoute(func(c *gin.Context) {
		response.FailNotFound(gin.H{"path": c.Request.URL.Path}, "接口不存在", c)
	})

	// 405 处理
	Router.NoMethod(func(c *gin.Context) {
		response.FailWithMessage("方法不被允许，请检查请求方式", c)
	})
}
