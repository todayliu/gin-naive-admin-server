package user

import (
	"github.com/gin-gonic/gin"
)

type _userService struct{}

var UserService = new(_userService)

func (us *_userService) GetUserInfoByUserId(c *gin.Context) {
	
}
