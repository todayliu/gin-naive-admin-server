package user

import "github.com/gin-gonic/gin"

type _userRouter struct{}

var UserRouter = new(_userRouter)

func (r *_userRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.POST("list", UserService.GetUserList)
		userRouter.GET("query/:id", UserService.QueryUser)
	}
}
