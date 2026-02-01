package login

import (
	"github.com/gin-gonic/gin"
)

var LoginRouter = new(_loginRouter)

type _loginRouter struct{}

func (l *_loginRouter) InitLoginRouter(Router *gin.RouterGroup) {
	loginRouter := Router.Group("login")
	{
		loginRouter.POST("captcha", LoginService.Captcha)
		loginRouter.POST("login", LoginService.Login)
	}
}
