package jwt_util

import (
	"gin-admin-server/global"

	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) (*CustomClaims, error) {
	token := c.Request.Header.Get("AccessToken")
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GNA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在AccessToken且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.ID
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.BaseClaims.ID
	}
}

// GetUserAccount 从Gin的Context中获取从jwt解析出来的用户账号
func GetUserAccount(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Account
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.Account
	}
}

// GetUserUuid 从Gin的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.UUID
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户信息
func GetUserInfo(c *gin.Context) *CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse
	}
}

// GetUserName 从Gin的Context中获取从jwt解析出来的用户名
func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.UName
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.UName
	}
}
