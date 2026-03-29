package position

import "github.com/gin-gonic/gin"

type _positionRouter struct{}

var PositionRouter = new(_positionRouter)

// InitPositionRouter 初始化职务管理模块路由
func (r *_positionRouter) InitPositionRouter(Router *gin.RouterGroup) {
	g := Router.Group("position")
	{
		g.GET("list", PositionService.GetPositionList)
		g.GET("export", PositionService.ExportPositions)
		g.GET("import-template", PositionService.DownloadPositionImportTemplate)
		g.POST("import", PositionService.ImportPositions)
		g.GET("query/:id", PositionService.QueryPosition)
		g.POST("add", PositionService.AddPosition)
		g.PUT("edit", PositionService.EditPosition)
		g.DELETE("delete/:id", PositionService.DeletePosition)
	}
}
