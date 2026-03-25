package profile

import "github.com/gin-gonic/gin"

type _profileRouter struct{}

var ProfileRouter = new(_profileRouter)

func (r *_profileRouter) InitProfileRouter(Router *gin.RouterGroup) {
	g := Router.Group("profile")
	{
		g.GET("info", ProfileService.GetInfo)
		g.PUT("info", ProfileService.UpdateInfo)
		g.PUT("password", ProfileService.UpdatePassword)
	}
}
