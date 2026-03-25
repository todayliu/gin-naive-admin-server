package permissionapi

import "github.com/gin-gonic/gin"

type _permissionAPIRouter struct{}

var PermissionAPIRouter = new(_permissionAPIRouter)

func (r *_permissionAPIRouter) InitPermissionAPIRouter(Router *gin.RouterGroup) {
	g := Router.Group("permission")
	{
		g.GET("button-codes", PermissionAPIService.ButtonCodes)
	}
}
