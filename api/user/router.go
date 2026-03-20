package user

import "github.com/gin-gonic/gin"

type _userRouter struct{}

var UserRouter = new(_userRouter)

// InitUserRouter 初始化用户模块路由
func (r *_userRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.GET("list", UserService.GetUserList)
		userRouter.GET("query/:id", UserService.QueryUser)
		userRouter.GET("roles/:id", UserService.GetUserRoles)
		userRouter.POST("add", UserService.AddUser)
		userRouter.PUT("edit", UserService.EditUser)
		userRouter.DELETE("delete/:id", UserService.DeleteUser)
	}
}
