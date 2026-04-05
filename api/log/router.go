package log

import "github.com/gin-gonic/gin"

type _logRouter struct{}

var LogRouter = new(_logRouter)

// InitLogRouter 日志模块路由
func (r *_logRouter) InitLogRouter(Router *gin.RouterGroup) {
	logRouter := Router.Group("log")
	{
		logRouter.GET("login/list", LogService.GetLoginLogList)
		logRouter.GET("oper/list", LogService.GetOperLogList)
	}
}
