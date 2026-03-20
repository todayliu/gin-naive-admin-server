package user

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		global.GNA_LOG.Error("获取用户总数失败" + err.Error())
		response.FailWithMessage("获取用户总数失败", c)
		return
	}
	limit := userRequest.PageSize
	offset := userRequest.PageSize * (userRequest.PageNo - 1)

	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&list).Error
	if err != nil {
		global.GNA_LOG.Error("获取用户列表失败：" + err.Error())
		response.FailWithMessage("获取用户列表失败", c)
		return
	}

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

func (us *_userService) AddUser(c *gin.Context) {
	var req UserAddRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, req)
		response.FailWithMessage(errMessage, c)
		return
	}

	var existUser SysUser
	err = global.GNA_DB.Where("account = ?", req.Account).First(&existUser).Error
	if err == nil {
		response.FailWithMessage("用户账号已存在", c)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GNA_LOG.Error("检查用户账号失败：" + err.Error())
		response.FailWithMessage("添加用户失败", c)
		return
	}

	user := SysUser{
		UUID:      utils.GenerateUUID(),
		Account:   req.Account,
		Password:  req.Password,
		UName:     req.UName,
		UNickname: req.UNickname,
		UMobile:   req.UMobile,
		UEmail:    req.UEmail,
		UAvatar:   req.UAvatar,
		Gender:    req.Gender,
		Status:    req.Status,
	}

	err = global.GNA_DB.Create(&user).Error
	if err != nil {
		global.GNA_LOG.Error("添加用户失败：" + err.Error())
		response.FailWithMessage("添加用户失败", c)
		return
	}

	response.Ok(c)
}

func (us *_userService) EditUser(c *gin.Context) {
	var req UserEditRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, req)
		response.FailWithMessage(errMessage, c)
		return
	}

	updates := map[string]interface{}{
		"account":    req.Account,
		"u_name":     req.UName,
		"u_nickname": req.UNickname,
		"u_mobile":   req.UMobile,
		"u_email":    req.UEmail,
		"u_avatar":   req.UAvatar,
		"gender":     req.Gender,
		"status":     req.Status,
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}

	err = global.GNA_DB.Model(&SysUser{}).Where("id = ?", req.ID).Updates(updates).Error
	if err != nil {
		global.GNA_LOG.Error("修改用户失败：" + err.Error())
		response.FailWithMessage("修改用户失败", c)
		return
	}

	response.Ok(c)
}

func (us *_userService) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}

	err := global.GNA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM sys_user_role WHERE sys_user_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("id = ?", id).Delete(&SysUser{}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		global.GNA_LOG.Error("删除用户失败：" + err.Error())
		response.FailWithMessage("删除用户失败", c)
		return
	}

	response.Ok(c)
}
