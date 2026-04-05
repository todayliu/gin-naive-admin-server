package middleware

import (
	"errors"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/jwt_util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("AccessToken")
		if accessToken == "" {
			response.FailWithMessageByToken("未登录或非法访问", c)
			c.Abort()
			return
		}

		j := jwt_util.NewJWT()
		claims, err := j.ParseToken(accessToken)
		if err != nil {
			if errors.Is(err, jwt_util.TokenExpired) {
				response.FailWithMessageByToken("授权已过期，请重新登录", c)
				c.Abort()
				return
			}
			response.FailWithMessageByToken(err.Error(), c)
			c.Abort()
			return
		}

		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			dr := time.Duration(global.GNA_CONFIG.Jwt.ExpiresTime) * time.Second
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			newToken, _ := j.CreateTokenByOldToken(accessToken, *claims)
			c.Header("New-Access-Token", newToken)
		}
		c.Set("claims", claims)
		c.Request = c.Request.WithContext(global.WithOperatorUserID(c.Request.Context(), claims.BaseClaims.ID))
		c.Next()
	}
}
