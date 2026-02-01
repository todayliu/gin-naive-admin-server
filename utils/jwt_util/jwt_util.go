package jwt_util

import (
	"errors"
	"gin-admin-server/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GNA_CONFIG.Jwt.SecretKey),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	bf := global.GNA_CONFIG.Jwt.BufferTime
	ep := global.GNA_CONFIG.Jwt.ExpiresTime
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: bf, // 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GNA"},                                             // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),                           // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ep) * time.Second)), // 过期时间 7天  配置文件
			Issuer:    global.GNA_CONFIG.Jwt.Issuer,                                        // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		// 1. 检查 Token 格式是否错误 (对应之前的 ValidationErrorMalformed)
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalformed
		}

		// 2. 检查 Token 是否已过期 (对应之前的 ValidationErrorExpired)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		}

		// 3. 检查 Token 是否尚未生效 (对应之前的 ValidationErrorNotValidYet)
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		}

		// 4. 检查签名是否无效 (比如被篡改了)
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, TokenInvalid
		}

		// 其他兜底错误
		return nil, TokenInvalid
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
