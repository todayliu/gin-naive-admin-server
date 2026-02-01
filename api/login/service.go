package login

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils"
	"gin-admin-server/utils/captcha_redis"
	"image/color"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = captcha_redis.NewDefaultRedisStore()
var LoginService = new(_LoginService)

type _LoginService struct{}

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
	}, "验证码获取成功", c)
}

func (ls *_LoginService) Login(c *gin.Context) {
	var loginRequest LoginRequest
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(loginRequest, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("登录成功", c)
}
