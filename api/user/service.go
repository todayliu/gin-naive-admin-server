package user

import (
	"crypto/sha256"
	"encoding/hex"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"
	"gin-admin-server/utils"
	"gin-admin-server/utils/validator"
	"strconv"
	"strings"

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
// @Summary     用户分页列表
// @Description 分页列表；筛选条件见 query 参数。
// @Tags        用户
// @Produce     json
// @Security    AccessToken
// @Param       pageNo query int true "页码"
// @Param       pageSize query int true "每页条数"
// @Param       username query string false "账号模糊"
// @Param       nickname query string false "昵称模糊"
// @Param       gender query string false "性别"
// @Param       status query string false "状态"
// @Param       departmentId query string false "部门 ID"
// @Success     200 {object} response.Response
// @Router      /user/list [get]
func (us *_userService) GetUserList(c *gin.Context) {
	var list []SysUser
	var userRequest UserPageRequest
	err := c.ShouldBindQuery(&userRequest)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, userRequest)
		response.FailWithMessage(errMessage, c)
		return
	}
	db := buildUserListQuery(c, &userRequest.UserListFilters)

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

	global.FillAuditDisplayNames(dbctx.Use(c), &list)

	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   userRequest.PageNo,
		PageSize: userRequest.PageSize,
	}, c)
}

// buildUserListQuery 列表/导出共用的筛选与数据范围（不含分页、排序）
func buildUserListQuery(c *gin.Context, filters *UserListFilters) *gorm.DB {
	db := dbctx.Use(c).Model(&SysUser{})

	if filters.Username != "" {
		db = db.Where("account LIKE ?", "%"+strings.TrimSpace(filters.Username)+"%")
	}
	if filters.UName != "" {
		db = db.Where("u_name LIKE ?", "%"+strings.TrimSpace(filters.UName)+"%")
	}
	nick := strings.TrimSpace(filters.Nickname)
	if nick == "" {
		nick = strings.TrimSpace(filters.UNickname)
	}
	if nick != "" {
		db = db.Where("u_nickname LIKE ?", "%"+nick+"%")
	}
	if filters.Gender != "" {
		if g, err := strconv.ParseUint(filters.Gender, 10, 32); err == nil {
			db = db.Where("gender = ?", uint(g))
		}
	}
	if filters.Status != "" {
		if s, err := strconv.ParseUint(filters.Status, 10, 32); err == nil {
			db = db.Where("status = ?", uint(s))
		}
	}
	if filters.DepartmentID != "" {
		if d, err := strconv.ParseUint(filters.DepartmentID, 10, 32); err == nil && d > 0 {
			did := uint(d)
			db = db.Where("department_id = ? OR EXISTS (SELECT 1 FROM sys_user_department sud WHERE sud.sys_user_id = sys_user.id AND sud.sys_department_id = ?)", did, did)
		}
	}

	return db
}

// QueryUser 查询用户详情（含角色ID列表）
// @Summary     用户详情
// @Tags        用户
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "用户 ID"
// @Success     200 {object} response.Response
// @Router      /user/query/{id} [get]
func (us *_userService) QueryUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var sysUser SysUser
	err := dbctx.Use(c).Where("id = ?", id).First(&sysUser).Error
	if err != nil {
		global.GNA_LOG.Error("用户查询失败：" + err.Error())
		response.FailWithMessage("用户查询失败", c)
		return
	}

	var roleIds []uint
	_ = dbctx.Use(c).Table("sys_user_role").Where("sys_user_id = ?", id).Pluck("sys_role_id", &roleIds)

	var departmentIds []uint
	_ = dbctx.Use(c).Model(&SysUserDepartment{}).Where("sys_user_id = ?", id).Order("sys_department_id").Pluck("sys_department_id", &departmentIds)
	if len(departmentIds) == 0 && sysUser.DepartmentId > 0 {
		departmentIds = []uint{sysUser.DepartmentId}
	}

	var positionIds []uint
	_ = dbctx.Use(c).Model(&SysUserJobLevel{}).Where("sys_user_id = ?", id).Order("sys_job_level_id").Pluck("sys_job_level_id", &positionIds)
	if len(positionIds) == 0 && sysUser.JobLevelID > 0 {
		positionIds = []uint{sysUser.JobLevelID}
	}

	global.FillAuditDisplayNames(dbctx.Use(c), &sysUser)

	resp := map[string]interface{}{
		"id":            sysUser.ID,
		"createTime":    sysUser.CreateTime,
		"updateTime":    sysUser.UpdateTime,
		"createBy":      sysUser.CreateBy,
		"updateBy":      sysUser.UpdateBy,
		"createByName":  sysUser.CreateByName,
		"updateByName":  sysUser.UpdateByName,
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
		"positionIds":   positionIds,
		"lastLoginTime": sysUser.LastLoginTime,
		"roleIds":       roleIds,
	}
	response.OkWithData(resp, c)
}

// GetUserRoles 获取用户所属角色ID列表
// @Summary     用户角色 ID 列表
// @Tags        用户
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "用户 ID"
// @Success     200 {object} response.Response
// @Router      /user/roles/{id} [get]
func (us *_userService) GetUserRoles(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var roleIds []uint
	result := dbctx.Use(c).Table("sys_user_role").Where("sys_user_id = ?", id).Pluck("sys_role_id", &roleIds)
	if result.Error != nil {
		global.GNA_LOG.Error("查询用户角色失败", zap.Error(result.Error))
		response.FailWithMessage("查询用户角色失败", c)
		return
	}
	response.OkWithData(map[string]interface{}{"roleIds": roleIds}, c)
}

// AddUser 新增用户
// @Summary     新增用户
// @Tags        用户
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body UserAddRequest true "用户"
// @Success     200 {object} response.Response
// @Router      /user/add [post]
func (us *_userService) AddUser(c *gin.Context) {
	var req UserAddRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, req)
		response.FailWithMessage(errMessage, c)
		return
	}

	var existUser SysUser
	err = dbctx.Use(c).Where("account = ?", req.Account).First(&existUser).Error
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

	primaryJobLevel := uint(0)
	if len(req.PositionIDs) > 0 {
		primaryJobLevel = req.PositionIDs[0]
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
		JobLevelID:   primaryJobLevel,
	}

	err = dbctx.Use(c).Create(&user).Error
	if err != nil {
		global.GNA_LOG.Error("添加用户失败：" + err.Error())
		response.FailWithMessage("添加用户失败", c)
		return
	}

	// 设置用户部门关联
	if len(deptIds) > 0 {
		for _, deptId := range deptIds {
			if err := dbctx.Use(c).Create(&SysUserDepartment{SysUserID: user.ID, SysDepartmentID: deptId}).Error; err != nil {
				global.GNA_LOG.Error("设置用户部门失败：" + err.Error())
				break
			}
		}
	}

	// 设置用户角色关联
	if len(req.RoleIds) > 0 {
		for _, roleId := range req.RoleIds {
			if err := dbctx.Use(c).Exec("INSERT INTO sys_user_role (sys_user_id, sys_role_id) VALUES (?, ?)", user.ID, roleId).Error; err != nil {
				global.GNA_LOG.Error("设置用户角色失败：" + err.Error())
				break
			}
		}
	}

	// 设置用户职务关联（多选）
	if len(req.PositionIDs) > 0 {
		for _, jlID := range req.PositionIDs {
			if err := dbctx.Use(c).Create(&SysUserJobLevel{SysUserID: user.ID, SysJobLevelID: jlID}).Error; err != nil {
				global.GNA_LOG.Error("设置用户职务失败：" + err.Error())
				break
			}
		}
	}

	response.Ok(c)
}

// EditUser 修改用户（密码可选，不传则不更新）
// @Summary     编辑用户
// @Tags        用户
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body UserEditRequest true "用户"
// @Success     200 {object} response.Response
// @Router      /user/edit [put]
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
	if req.PositionIDs != nil {
		jl := uint(0)
		if len(req.PositionIDs) > 0 {
			jl = req.PositionIDs[0]
		}
		updates["job_level_id"] = jl
	}
	if req.Password != "" {
		updates["password"] = hashPasswordForStorage(req.Password)
	}

	err = dbctx.Use(c).Model(&SysUser{}).Where("id = ?", req.ID).Updates(updates).Error
	if err != nil {
		global.GNA_LOG.Error("修改用户失败：" + err.Error())
		response.FailWithMessage("修改用户失败", c)
		return
	}

	// 更新用户部门关联
	if req.DepartmentIds != nil || req.DepartmentId > 0 {
		if err := dbctx.Use(c).Where("sys_user_id = ?", req.ID).Delete(&SysUserDepartment{}).Error; err != nil {
			global.GNA_LOG.Error("清除用户部门失败：" + err.Error())
		}
		for _, deptId := range deptIds {
			if err := dbctx.Use(c).Create(&SysUserDepartment{SysUserID: req.ID, SysDepartmentID: deptId}).Error; err != nil {
				global.GNA_LOG.Error("设置用户部门失败：" + err.Error())
				break
			}
		}
	}

	// 更新用户角色关联
	if req.RoleIds != nil {
		if err := dbctx.Use(c).Exec("DELETE FROM sys_user_role WHERE sys_user_id = ?", req.ID).Error; err != nil {
			global.GNA_LOG.Error("清除用户角色失败：" + err.Error())
			response.FailWithMessage("修改用户失败", c)
			return
		}
		for _, roleId := range req.RoleIds {
			if err := dbctx.Use(c).Exec("INSERT INTO sys_user_role (sys_user_id, sys_role_id) VALUES (?, ?)", req.ID, roleId).Error; err != nil {
				global.GNA_LOG.Error("设置用户角色失败：" + err.Error())
				response.FailWithMessage("修改用户失败", c)
				return
			}
		}
	}

	// 更新用户职务关联（多选）
	if req.PositionIDs != nil {
		if err := dbctx.Use(c).Where("sys_user_id = ?", req.ID).Delete(&SysUserJobLevel{}).Error; err != nil {
			global.GNA_LOG.Error("清除用户职务失败：" + err.Error())
		}
		for _, jlID := range req.PositionIDs {
			if err := dbctx.Use(c).Create(&SysUserJobLevel{SysUserID: req.ID, SysJobLevelID: jlID}).Error; err != nil {
				global.GNA_LOG.Error("设置用户职务失败：" + err.Error())
				break
			}
		}
	}

	response.Ok(c)
}

// DeleteUser 删除用户（永久删除，同时删除用户角色关联）
// @Summary     删除用户
// @Tags        用户
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "用户 ID"
// @Success     200 {object} response.Response
// @Router      /user/delete/{id} [delete]
func (us *_userService) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}

	err := dbctx.Use(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM sys_user_role WHERE sys_user_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Where("sys_user_id = ?", id).Delete(&SysUserDepartment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("sys_user_id = ?", id).Delete(&SysUserJobLevel{}).Error; err != nil {
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
