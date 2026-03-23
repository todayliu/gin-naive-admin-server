package user

import (
	"crypto/sha256"
	"encoding/hex"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type _userService struct{}

var UserService = new(_userService)

// hashPasswordForStorage 对明文密码做 SHA256 哈希后存储，与登录校验逻辑一致（登录时前端发送 SHA256(明文)）
func hashPasswordForStorage(plainPassword string) string {
	hash := sha256.Sum256([]byte(plainPassword))
	return hex.EncodeToString(hash[:])
}

// GetUserList 分页查询用户列表
func (us *_userService) GetUserList(c *gin.Context) {
	var list []SysUser
	var userRequest UserPageRequest
	err := c.ShouldBindQuery(&userRequest)
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

// QueryUser 查询用户详情（含角色ID列表）
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

	var roleIds []uint
	_ = global.GNA_DB.Table("sys_user_role").Where("sys_user_id = ?", id).Pluck("sys_role_id", &roleIds)

	var departmentIds []uint
	_ = global.GNA_DB.Model(&SysUserDepartment{}).Where("sys_user_id = ?", id).Order("sys_department_id").Pluck("sys_department_id", &departmentIds)
	if len(departmentIds) == 0 && sysUser.DepartmentId > 0 {
		departmentIds = []uint{sysUser.DepartmentId}
	}

	resp := map[string]interface{}{
		"id":            sysUser.ID,
		"createTime":    sysUser.CreateTime,
		"updateTime":    sysUser.UpdateTime,
		"account":       sysUser.Account,
		"uName":         sysUser.UName,
		"uNickname":     sysUser.UNickname,
		"uMobile":       sysUser.UMobile,
		"uEmail":        sysUser.UEmail,
		"uAvatar":       sysUser.UAvatar,
		"gender":        sysUser.Gender,
		"status":        sysUser.Status,
		"departmentId":  sysUser.DepartmentId,
		"departmentIds": departmentIds,
		"lastLoginTime": sysUser.LastLoginTime,
		"roleIds":       roleIds,
	}
	response.OkWithData(resp, c)
}

// GetUserRoles 获取用户所属角色ID列表
func (us *_userService) GetUserRoles(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var roleIds []uint
	result := global.GNA_DB.Table("sys_user_role").Where("sys_user_id = ?", id).Pluck("sys_role_id", &roleIds)
	if result.Error != nil {
		global.GNA_LOG.Error("查询用户角色失败", zap.Error(result.Error))
		response.FailWithMessage("查询用户角色失败", c)
		return
	}
	response.OkWithData(map[string]interface{}{"roleIds": roleIds}, c)
}

// AddUser 新增用户
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

	deptIds := req.DepartmentIds
	if len(deptIds) == 0 && req.DepartmentId > 0 {
		deptIds = []uint{req.DepartmentId}
	}
	primaryDeptId := uint(0)
	if len(deptIds) > 0 {
		primaryDeptId = deptIds[0]
	}

	hashedPassword := hashPasswordForStorage(req.Password)
	user := SysUser{
		UUID:         utils.GenerateUUID(),
		Account:      req.Account,
		Password:     hashedPassword,
		UName:        req.UName,
		UNickname:    req.UNickname,
		UMobile:      req.UMobile,
		UEmail:       req.UEmail,
		UAvatar:      req.UAvatar,
		Gender:       req.Gender,
		Status:       req.Status,
		DepartmentId: primaryDeptId,
	}

	err = global.GNA_DB.Create(&user).Error
	if err != nil {
		global.GNA_LOG.Error("添加用户失败：" + err.Error())
		response.FailWithMessage("添加用户失败", c)
		return
	}

	// 设置用户部门关联
	if len(deptIds) > 0 {
		for _, deptId := range deptIds {
			if err := global.GNA_DB.Create(&SysUserDepartment{SysUserID: user.ID, SysDepartmentID: deptId}).Error; err != nil {
				global.GNA_LOG.Error("设置用户部门失败：" + err.Error())
				break
			}
		}
	}

	// 设置用户角色关联
	if len(req.RoleIds) > 0 {
		for _, roleId := range req.RoleIds {
			if err := global.GNA_DB.Exec("INSERT INTO sys_user_role (sys_user_id, sys_role_id) VALUES (?, ?)", user.ID, roleId).Error; err != nil {
				global.GNA_LOG.Error("设置用户角色失败：" + err.Error())
				break
			}
		}
	}

	response.Ok(c)
}

// EditUser 修改用户（密码可选，不传则不更新）
func (us *_userService) EditUser(c *gin.Context) {
	var req UserEditRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, req)
		response.FailWithMessage(errMessage, c)
		return
	}

	deptIds := req.DepartmentIds
	if len(deptIds) == 0 && req.DepartmentId > 0 {
		deptIds = []uint{req.DepartmentId}
	}
	primaryDeptId := uint(0)
	if len(deptIds) > 0 {
		primaryDeptId = deptIds[0]
	}

	updates := map[string]interface{}{
		"account":       req.Account,
		"u_name":        req.UName,
		"u_nickname":    req.UNickname,
		"u_mobile":      req.UMobile,
		"u_email":       req.UEmail,
		"u_avatar":      req.UAvatar,
		"gender":        req.Gender,
		"status":        req.Status,
		"department_id": primaryDeptId,
	}
	if req.Password != "" {
		updates["password"] = hashPasswordForStorage(req.Password)
	}

	err = global.GNA_DB.Model(&SysUser{}).Where("id = ?", req.ID).Updates(updates).Error
	if err != nil {
		global.GNA_LOG.Error("修改用户失败：" + err.Error())
		response.FailWithMessage("修改用户失败", c)
		return
	}

	// 更新用户部门关联
	if req.DepartmentIds != nil || req.DepartmentId > 0 {
		if err := global.GNA_DB.Where("sys_user_id = ?", req.ID).Delete(&SysUserDepartment{}).Error; err != nil {
			global.GNA_LOG.Error("清除用户部门失败：" + err.Error())
		}
		for _, deptId := range deptIds {
			if err := global.GNA_DB.Create(&SysUserDepartment{SysUserID: req.ID, SysDepartmentID: deptId}).Error; err != nil {
				global.GNA_LOG.Error("设置用户部门失败：" + err.Error())
				break
			}
		}
	}

	// 更新用户角色关联
	if req.RoleIds != nil {
		if err := global.GNA_DB.Exec("DELETE FROM sys_user_role WHERE sys_user_id = ?", req.ID).Error; err != nil {
			global.GNA_LOG.Error("清除用户角色失败：" + err.Error())
			response.FailWithMessage("修改用户失败", c)
			return
		}
		for _, roleId := range req.RoleIds {
			if err := global.GNA_DB.Exec("INSERT INTO sys_user_role (sys_user_id, sys_role_id) VALUES (?, ?)", req.ID, roleId).Error; err != nil {
				global.GNA_LOG.Error("设置用户角色失败：" + err.Error())
				response.FailWithMessage("修改用户失败", c)
				return
			}
		}
	}

	response.Ok(c)
}

// DeleteUser 删除用户（永久删除，同时删除用户角色关联）
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
		if err := tx.Where("sys_user_id = ?", id).Delete(&SysUserDepartment{}).Error; err != nil {
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
