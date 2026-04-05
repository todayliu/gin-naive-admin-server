package login

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"gin-admin-server/api/log"
	"gin-admin-server/api/user"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils"
	"gin-admin-server/utils/captcha_redis"
	"gin-admin-server/utils/jwt_util"
	"gin-admin-server/utils/validator"
	"image/color"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var store = captcha_redis.NewDefaultRedisStore()
var LoginService = new(_LoginService)

type _LoginService struct{}

// Captcha 获取图形验证码
// @Summary     获取验证码
// @Description 返回 base64 图片与 captchaId、盐；登录前调用。验证码答案存于 Redis，键为 CAPTCHA_{captchaId}。
// @Tags        认证
// @Produce     json
// @Success     200 {object} response.Response
// @Router      /login/captcha [post]
func (ls *_LoginService) Captcha(c *gin.Context) {
	driver := base64Captcha.DriverString{
		Width:           global.GNA_CONFIG.Captcha.ImgWidth,
		Height:          global.GNA_CONFIG.Captcha.ImgHeight,
		NoiseCount:      global.GNA_CONFIG.Captcha.NoiseCount,
		ShowLineOptions: global.GNA_CONFIG.Captcha.ShowLineOptions,
		Length:          global.GNA_CONFIG.Captcha.KeyLong,
		Source:          global.GNA_CONFIG.Captcha.Source,
		BgColor:         &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	cp := base64Captcha.NewCaptcha(&driver, store.UseWithCtx(c))
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.GNA_LOG.Error("验证码获取失败!", zap.Error(err))
		response.FailWithMessage("验证码获取失败", c)
		return
	}

	response.OkWithDetailed(CaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.GNA_CONFIG.Captcha.KeyLong,
		Salt:          utils.GenerateRandomString(24),
	}, "验证码获取成功", c)
}

// Login 用户登录
// @Summary     登录
// @Description password 为前端对明文密码做 SHA256 后，再与 salt、timestamp 组合 SHA256 的十六进制字符串；需携带验证码与 nonce（防重放）。
// @Tags        认证
// @Accept      json
// @Produce     json
// @Param       body body LoginRequest true "登录参数"
// @Success     200 {object} response.Response
// @Router      /login/login [post]
func (ls *_LoginService) Login(c *gin.Context) {
	var loginRequest LoginRequest
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, loginRequest)
		global.GNA_LOG.Error(errMessage)
		response.FailWithMessage(errMessage, c)
		return
	}

	redisNonce := global.GNA_REDIS.Get(c.Request.Context(), "login:nonce:"+loginRequest.Username+":"+loginRequest.Nonce)
	_, err = redisNonce.Result()
	if err == nil {
		global.GNA_LOG.Warn("请求过于频繁或重复登录: " + loginRequest.Username)
		response.FailWithMessage("请勿重复操作", c)
		return
	}
	if !errors.Is(err, redis.Nil) {
		global.GNA_LOG.Error("查询异常: " + err.Error())
		response.FailWithMessage("服务器繁忙，请稍后再试", c)
		return
	}
	global.GNA_REDIS.Set(c.Request.Context(), "login:nonce:"+loginRequest.Username+":"+loginRequest.Nonce, loginRequest.Nonce, 5*time.Minute)

	// 先校验验证码，减少暴力破解风险
	smsCode, err := global.GNA_REDIS.Get(c.Request.Context(), "CAPTCHA_"+loginRequest.CaptchaId).Result()
	if err != nil || strings.ToLower(smsCode) != strings.ToLower(loginRequest.Captcha) {
		global.GNA_LOG.Error("验证码错误")
		log.SaveLoginLogAsync(0, loginRequest.Username, c.ClientIP(), 2, "验证码错误")
		response.FailWithMessage("验证码错误，请重试", c)
		return
	}

	var userInfo user.SysUser
	err = global.GNA_DB.Where("account = ?", loginRequest.Username).First(&userInfo).Error
	if err != nil {
		global.GNA_LOG.Error("用户不存在或查询失败", zap.Error(err))
		log.SaveLoginLogAsync(0, loginRequest.Username, c.ClientIP(), 2, "用户不存在")
		response.FailWithMessage("用户名或密码错误，请重试", c)
		return
	}

	// 判断密码
	newSalt := global.GNA_CONFIG.System.Name + loginRequest.Salt
	combined := fmt.Sprintf("%s:%s:%d", userInfo.Password, newSalt, loginRequest.Timestamp)
	hash := sha256.Sum256([]byte(combined))
	expectedHash := hex.EncodeToString(hash[:])

	if expectedHash != loginRequest.Password {
		global.GNA_LOG.Error("密码错误")
		log.SaveLoginLogAsync(userInfo.ID, userInfo.Account, c.ClientIP(), 2, "密码错误")
		response.FailWithMessage("用户名或密码错误，请重试", c)
		return
	}

	//签发 Token
	ls.CreateToken(c, userInfo)
}

func (ls *_LoginService) CreateToken(c *gin.Context, userInfo user.SysUser) {
	j := &jwt_util.JWT{SigningKey: []byte(global.GNA_CONFIG.Jwt.SecretKey)} // 唯一签名
	claims := j.CreateClaims(jwt_util.BaseClaims{
		UUID:    userInfo.UUID,
		ID:      userInfo.ID,
		Account: userInfo.Account,
		UName:   userInfo.UName,
		UMobile: userInfo.UMobile,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GNA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败!", c)
		return
	}

	db := global.GNA_DB.WithContext(global.WithOperatorUserID(context.Background(), userInfo.ID))
	_ = db.Model(&user.SysUser{}).Where("id = ?", userInfo.ID).Update("last_login_time", time.Now()).Error
	log.SaveLoginLogAsync(userInfo.ID, userInfo.Account, c.ClientIP(), 1, "登录成功")

	roles := loadUserRoleCodes(userInfo.ID)

	response.OkWithDetailed(LoginResponse{
		UserInfo:  userInfo,
		Token:     token,
		Roles:     roles,
		Codes:     nil,
		ExpiresIn: global.GNA_CONFIG.Jwt.ExpiresTime,
	}, "登录成功", c)
}
