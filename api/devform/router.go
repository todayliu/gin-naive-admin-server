package devform

import "github.com/gin-gonic/gin"

type _devformRouter struct{}

var DevformRouter = new(_devformRouter)

func (r *_devformRouter) InitDevformRouter(Router *gin.RouterGroup) {
	g := Router.Group("devform")
	{
		g.GET("list", DevformService.GetFormList)
		g.GET("query/:id", DevformService.QueryForm)
		g.POST("add", DevformService.AddForm)
		g.PUT("edit", DevformService.EditForm)
		g.DELETE("delete/:id", DevformService.DeleteForm)
		g.POST("field/save", DevformService.SaveField)
		g.DELETE("field/delete/:id", DevformService.DeleteField)
		g.POST("sync/:id", DevformService.SyncDB)
		g.GET("download/:id", DevformService.DownloadCode)
	}
}
