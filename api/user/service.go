package user

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
)

type _userService struct{}

var UserService = new(_userService)

func (us *_userService) GetUserList(c *gin.Context) {
	var list []SysUser
	var userRequest UserPageRequest
	err := c.ShouldBindJSON(&userRequest)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, userRequest)
		response.FailWithMessage(errMessage, c)
		return
	}
	db := global.GNA_DB.Model(&SysUser{})

	var total int64

	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("获取角色总数失败" + err.Error())
		response.FailWithMessage("获取角色总数失败", c)
		return
	}
	limit := userRequest.PageSize
	offset := userRequest.PageSize * (userRequest.PageNo - 1)

	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&list).Error

	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   userRequest.PageNo,
		PageSize: userRequest.PageSize,
	}, c)
}

func (us *_userService) QueryUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var sysUser SysUser
	err := global.GNA_DB.Where("id = ?", id).First(&sysUser).Error
	if err != nil {
		global.GNA_LOG.Error("用户查询失败：" + err.Error())
		response.FailWithMessage("用户查询失败", c)
		return
	}

	response.OkWithData(sysUser, c)
}
