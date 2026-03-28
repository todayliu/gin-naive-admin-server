package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck 服务健康检查。
// @Summary     健康检查
// @Description 用于探活或负载均衡，无需鉴权。
// @Tags        系统
// @Produce     json
// @Success     200 {string} string "服务状态正常"
// @Router      /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "服务状态正常")
}
