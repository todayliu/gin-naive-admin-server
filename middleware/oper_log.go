package middleware

import (
	"gin-admin-server/api/log"
	"gin-admin-server/global"
	"gin-admin-server/utils/jwt_util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OperLog 记录变更类请求的操作日志（异步）
func OperLog() gin.HandlerFunc {
	prefix := global.GNA_CONFIG.Router.RouterPrefix
	if prefix == "" {
		prefix = "/api"
	}
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		method := c.Request.Method
		if method == "GET" || method == "HEAD" {
			c.Next()
			return
		}
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		if strings.HasPrefix(path, prefix+"/log/") {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		uid := jwt_util.GetUserID(c)
		acc := jwt_util.GetUserAccount(c)
		title := operTitle(path)
		log.SaveOperLogAsync(title, method, path, uid, acc, c.ClientIP(), latency, c.Writer.Status())
	}
}

func operTitle(path string) string {
	if strings.Contains(path, "/user/") {
		return "用户管理"
	}
	if strings.Contains(path, "/role/") {
		return "角色管理"
	}
	if strings.Contains(path, "/menu/") {
		return "菜单管理"
	}
	if strings.Contains(path, "/department/") {
		return "部门管理"
	}
	if strings.Contains(path, "/dict/") {
		return "字典管理"
	}
	if strings.Contains(path, "/position/") {
		return "职务管理"
	}
	if strings.Contains(path, "/profile/") {
		return "个人中心"
	}
	if strings.Contains(path, "/file/") {
		return "文件上传"
	}
	if strings.Contains(path, "/config/") {
		return "参数配置"
	}
	return "系统"
}
