package middleware

import (
	"fmt"
	"gin-admin-server/config"
	"gin-admin-server/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	ruleMode := global.GNA_CONFIG.Cors.Mode
	switch ruleMode {
	case "whitelist":
		return CorsByWhitelist()
	default:
		return CorsByAll()
	}
}

func CorsByAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,UUID")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Refresh-Token, Refresh-Token-Expires-Time, New-Access-Token")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func CorsByWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			// 同源请求（如直接访问 API、健康检查）无 Origin 头，放行
			host := c.Request.Host
			allowedHosts := []string{
				fmt.Sprintf("127.0.0.1:%d", global.GNA_CONFIG.System.Port),
				fmt.Sprintf("localhost:%d", global.GNA_CONFIG.System.Port),
			}
			for _, h := range allowedHosts {
				if host == h {
					c.Next()
					return
				}
			}
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		whitelist := checkCors(origin)

		if whitelist == nil {
			if c.Request.Method == "GET" && (c.Request.URL.Path == "/api/health" || c.Request.URL.Path == "/health") {
				c.Next()
				return
			}
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 通过检查, 添加请求头
		c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)
		c.Header("Access-Control-Allow-Headers", whitelist.AllowHeaders)
		c.Header("Access-Control-Allow-Methods", whitelist.AllowMethods)
		c.Header("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
		if whitelist.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 处理请求
		c.Next()
	}
}

func checkCors(currentOrigin string) *config.CORSWhitelist {
	for _, whitelist := range global.GNA_CONFIG.Cors.Whitelist {
		// 遍历配置中的跨域头，寻找匹配项
		if currentOrigin == whitelist.AllowOrigin {
			return &whitelist
		}
	}
	return nil
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
