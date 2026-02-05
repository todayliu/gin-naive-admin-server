package menu

import "github.com/gin-gonic/gin"

type _menuRouter struct{}

var MenuRouter = new(_menuRouter)

func (l *_menuRouter) InitMenuRouter(Router *gin.RouterGroup) {
	menuRouter := Router.Group("menu")
	{
		menuRouter.POST("router", MenuService.InitMenuList)
		menuRouter.POST("list", MenuService.GetAllMenuList)
		menuRouter.POST("edit", MenuService.UpdateMenu)
		menuRouter.DELETE("delete/:id", MenuService.DeleteMenu)
	}
}
