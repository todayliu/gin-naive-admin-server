package login

import (
	"gin-admin-server/model/response"

	"github.com/gin-gonic/gin"
)

var LoginService = new(_LoginService)

type _LoginService struct{}

func (s *_LoginService) Login(c *gin.Context) {
	response.OkWithMessage("请求成功", c)
}
