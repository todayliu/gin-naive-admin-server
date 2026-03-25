package middleware

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/permission"
	"gin-admin-server/utils/jwt_util"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIPermission 基于权限码的接口鉴权（在 JWT 之后注册）
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
		if permission.IsSuperUser(uid) {
			c.Next()
			return
		}
		fullPath := c.FullPath()
		method := c.Request.Method
		if authOnlyWhitelist(method, fullPath) {
			c.Next()
			return
		}
		code, need := permission.RequiredCode(method, fullPath)
		if !need {
			c.Next()
			return
		}
		if permission.UserHasPermCode(uid, code) {
			c.Next()
			return
		}
		response.FailWithMessage("无访问权限: "+code, c)
		c.Abort()
	}
}

const httpMethodOptions = "OPTIONS"

func authOnlyWhitelist(method, fullPath string) bool {
	p := global.GNA_CONFIG.Router.RouterPrefix
	if p == "" {
		p = "/api"
	}
	if method == "GET" && fullPath == p+"/menu/router" {
		return true
	}
	if strings.HasPrefix(fullPath, p+"/profile/") {
		return true
	}
	// 字典按类型编码查询（下拉）
	if method == "GET" && strings.Contains(fullPath, "/dict/data/") {
		return true
	}
	return false
}
