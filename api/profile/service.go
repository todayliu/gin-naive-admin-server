package profile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gin-admin-server/api/user"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"
	"gin-admin-server/utils"
	"gin-admin-server/utils/jwt_util"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type _profileService struct{}

var ProfileService = new(_profileService)

// GetInfo 当前登录用户资料
// @Summary     个人资料
// @Tags        个人中心
// @Produce     json
// @Security    AccessToken
// @Success     200 {object} response.Response
// @Router      /profile/info [get]
func (s *_profileService) GetInfo(c *gin.Context) {
	uid := jwt_util.GetUserID(c)
	if uid == 0 {
		response.FailWithMessage("未登录", c)
		return
	}
	var u user.SysUser
	if err := dbctx.Use(c).Where("id = ?", uid).First(&u).Error; err != nil {
		response.FailWithMessage("用户不存在", c)
		return
	}
	roleCodes, _ := user.LoadUserRoleCodes(uid)
	response.OkWithData(map[string]interface{}{
		"id":            u.ID,
		"createTime":    u.CreateTime,
		"updateTime":    u.UpdateTime,
		"account":       u.Account,
		"uName":         u.UName,
		"uNickname":     u.UNickname,
		"uMobile":       u.UMobile,
		"uEmail":        u.UEmail,
		"uAvatar":       u.UAvatar,
		"gender":        u.Gender,
		"status":        u.Status,
		"departmentId":  u.DepartmentId,
		"lastLoginTime": u.LastLoginTime,
		"cryptoSalt":    utils.GenerateRandomString(24),
		"roles":         roleCodes,
	}, c)
}

// UpdateInfo 更新个人资料
// @Summary     更新个人资料
// @Tags        个人中心
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body ProfileUpdateRequest true "请求体"
// @Success     200 {object} response.Response
// @Router      /profile/info [put]
func (s *_profileService) UpdateInfo(c *gin.Context) {
	uid := jwt_util.GetUserID(c)
	if uid == 0 {
		response.FailWithMessage("未登录", c)
		return
	}
	var req ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	err := dbctx.Use(c).Model(&user.SysUser{}).Where("id = ?", uid).Updates(map[string]interface{}{
		"u_nickname": req.UNickname,
		"u_mobile":   req.UMobile,
		"u_email":    req.UEmail,
		"u_avatar":   req.UAvatar,
	}).Error
	if err != nil {
		global.GNA_LOG.Error("更新资料失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.Ok(c)
}

// UpdatePassword 修改登录密码
// @Summary     修改密码
// @Description 与登录相同：传输形态为二次 SHA256；需 salt、timestamp。
// @Tags        个人中心
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body ProfilePasswordRequest true "请求体"
// @Success     200 {object} response.Response
// @Router      /profile/password [put]
func (s *_profileService) UpdatePassword(c *gin.Context) {
	uid := jwt_util.GetUserID(c)
	if uid == 0 {
		response.FailWithMessage("未登录", c)
		return
	}
	var req ProfilePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	var u user.SysUser
	if err := dbctx.Use(c).Where("id = ?", uid).First(&u).Error; err != nil {
		response.FailWithMessage("用户不存在", c)
		return
	}
	newSalt := global.GNA_CONFIG.System.Name + req.Salt
	combinedOld := fmt.Sprintf("%s:%s:%d", u.Password, newSalt, req.Timestamp)
	sumOld := sha256.Sum256([]byte(combinedOld))
	hashOld := hex.EncodeToString(sumOld[:])
	if hashOld != req.OldPassword {
		response.FailWithMessage("原密码错误", c)
		return
	}
	newStored := hashStoredPassword(req.NewPassword)
	if err := dbctx.Use(c).Model(&user.SysUser{}).Where("id = ?", uid).Update("password", newStored).Error; err != nil {
		response.FailWithMessage("修改密码失败", c)
		return
	}
	response.Ok(c)
}

func hashStoredPassword(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
