package middleware

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/jwt_util"

	"github.com/gin-gonic/gin"
)

const httpMethodOptions = "OPTIONS"

// APIPermission 在 JWT 之后注册：未开启 relax-api-auth 时要求已登录；不按权限码拦截接口。
func APIPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == httpMethodOptions {
			c.Next()
			return
		}
		if global.GNA_CONFIG.Security.RelaxApiAuth {
			c.Next()
			return
		}
		uid := jwt_util.GetUserID(c)
		if uid == 0 {
			response.FailWithMessageByToken("未登录或非法访问", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
