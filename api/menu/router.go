package menu

import "github.com/gin-gonic/gin"

type _menuRouter struct{}

var MenuRouter = new(_menuRouter)

// InitMenuRouter 初始化菜单模块路由
func (l *_menuRouter) InitMenuRouter(Router *gin.RouterGroup) {
	menuRouter := Router.Group("menu")
	{
		menuRouter.GET("router", MenuService.InitMenuList)
		menuRouter.GET("list", MenuService.GetAllMenuList)
		menuRouter.PUT("edit", MenuService.UpdateMenu)
		menuRouter.POST("add", MenuService.AddMenu)
		menuRouter.DELETE("delete/:id", MenuService.DeleteMenu)
	}
}
