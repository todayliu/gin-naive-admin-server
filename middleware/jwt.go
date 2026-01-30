package middleware

import (
	"gin-admin-server/model/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Access-Token")
		if accessToken == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}

	}

	//if jwtService.IsBlacklist(token) {
	//	response.FailWithDetailed(gin.H{"reload": true}, "您的帐户异地登陆或令牌失效", c)
	//	c.Abort()
	//	return
	//}
}
