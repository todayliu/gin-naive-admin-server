package department

import "github.com/gin-gonic/gin"

type _departmentRouter struct{}

var DepartmentRouter = new(_departmentRouter)

// InitDepartmentRouter 初始化部门模块路由
func (r *_departmentRouter) InitDepartmentRouter(Router *gin.RouterGroup) {
	departmentRouter := Router.Group("department")
	{
		departmentRouter.GET("list", DepartmentService.GetDepartmentList)
		departmentRouter.GET("export", DepartmentService.ExportDepartmentPaths)
		departmentRouter.PUT("edit", DepartmentService.UpdateDepartment)
		departmentRouter.DELETE("delete/:id", DepartmentService.DeleteDepartment)
	}
}
