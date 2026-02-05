package role

import "github.com/gin-gonic/gin"

type _roleRouter struct{}

var RoleRouter = new(_roleRouter)

func (r *_roleRouter) InitRoleRouter(Router *gin.RouterGroup) {
	roleRouter := Router.Group("role")
	{
		roleRouter.POST("list", RoleService.GetRoleList)
		roleRouter.POST("add", RoleService.AddRole)
		roleRouter.GET("query/:id", RoleService.QueryRole)
		roleRouter.POST("edit", RoleService.EditRole)
		roleRouter.DELETE("delete/:id", RoleService.DeleteRole)
		roleRouter.GET("powerTree/:id", RoleService.GetPowerTree)
		roleRouter.POST("setRolePower", RoleService.SetRolePower)
	}
}
