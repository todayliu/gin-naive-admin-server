package dashboard

import "github.com/gin-gonic/gin"

type _dashboardRouter struct{}

var DashboardRouter = new(_dashboardRouter)

func (r *_dashboardRouter) InitDashboardRouter(Router *gin.RouterGroup) {
	Router.GET("dashboard/stats", DashboardService.Stats)
}
