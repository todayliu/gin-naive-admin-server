package role

import (
	"gin-admin-server/api/menu"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type _roleService struct{}

var RoleService = new(_roleService)

// GetRoleList 分页查询角色列表
// @Summary     角色分页列表
// @Tags        角色
// @Produce     json
// @Security    AccessToken
// @Param       pageNo query int true "页码"
// @Param       pageSize query int true "每页条数"
// @Param       name query string false "角色名模糊"
// @Param       code query string false "角色编码模糊"
// @Success     200 {object} response.Response
// @Router      /role/list [get]
func (r *_roleService) GetRoleList(c *gin.Context) {
	var list []SysRole
	var roleRequest RolePageRequest
	err := c.ShouldBindQuery(&roleRequest)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, roleRequest)
		response.FailWithMessage(errMessage, c)
		return
	}
	db := dbctx.Use(c).Model(&SysRole{})
	if roleRequest.Name != "" {
		db = db.Where("name LIKE ?", "%"+roleRequest.Name+"%")
	}
	if roleRequest.Code != "" {
		db = db.Where("code LIKE ?", roleRequest.Code)
	}
	if roleRequest.Status != nil {
		db = db.Where("status = ?", *roleRequest.Status)
	}
	var total int64

	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("获取角色总数失败" + err.Error())
		response.FailWithMessage("获取角色总数失败", c)
		return
	}
	limit := roleRequest.PageSize
	offset := roleRequest.PageSize * (roleRequest.PageNo - 1)

	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&list).Error
	if err != nil {
		global.GNA_LOG.Error("获取角色列表失败：" + err.Error())
		response.FailWithMessage("获取角色列表失败", c)
		return
	}
	global.FillAuditDisplayNames(dbctx.Use(c), &list)

	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   roleRequest.PageNo,
		PageSize: roleRequest.PageSize,
	}, c)
}

// AddRole 新增角色
// @Summary     新增角色
// @Tags        角色
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysRole true "角色"
// @Success     200 {object} response.Response
// @Router      /role/add [post]
func (r *_roleService) AddRole(c *gin.Context) {
	var roleReq SysRole
	err := c.ShouldBindJSON(&roleReq)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, roleReq)
		response.FailWithMessage(errMessage, c)
		return
	}

	err = dbctx.Use(c).Create(&roleReq).Error
	if err != nil {
		global.GNA_LOG.Error("添加角色失败：" + err.Error())
		response.FailWithMessage("添加角色失败", c)
		return
	}

	response.Ok(c)
}

// QueryRole 查询角色详情
// @Summary     角色详情
// @Tags        角色
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "角色 ID"
// @Success     200 {object} response.Response
// @Router      /role/query/{id} [get]
func (r *_roleService) QueryRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var sysRole SysRole
	err := dbctx.Use(c).Where("id = ?", id).First(&sysRole).Error
	if err != nil {
		global.GNA_LOG.Error("查询角色失败：" + err.Error())
		response.FailWithMessage("查询角色失败", c)
		return
	}

	global.FillAuditDisplayNames(dbctx.Use(c), &sysRole)
	response.OkWithData(sysRole, c)
}

// EditRole 修改角色
// @Summary     编辑角色
// @Tags        角色
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysRole true "角色（含 id）"
// @Success     200 {object} response.Response
// @Router      /role/edit [put]
func (r *_roleService) EditRole(c *gin.Context) {
	var roleReq SysRole
	err := c.ShouldBindJSON(&roleReq)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, roleReq)
		response.FailWithMessage(errMessage, c)
		return
	}

	err = dbctx.Use(c).Model(&roleReq).Updates(roleReq).Error
	if err != nil {
		global.GNA_LOG.Error("修改角色失败：" + err.Error())
		response.FailWithMessage("修改角色失败", c)
		return
	}

	response.Ok(c)
}

// DeleteRole 删除角色（永久删除，同时删除用户角色、角色菜单关联）
// @Summary     删除角色
// @Tags        角色
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "角色 ID"
// @Success     200 {object} response.Response
// @Router      /role/delete/{id} [delete]
func (r *_roleService) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}

	err := dbctx.Use(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM sys_user_role WHERE sys_role_id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM sys_role_menu WHERE sys_role_id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Where("id = ?", id).Delete(&SysRole{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.GNA_LOG.Error("删除角色失败：" + err.Error())
		response.FailWithMessage("删除角色失败", c)
		return
	}

	response.Ok(c)
}

// GetPowerTree 获取角色权限树（全部菜单树 + 角色已选菜单ID）
// @Summary     角色权限树
// @Tags        角色
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "角色 ID"
// @Success     200 {object} response.Response
// @Router      /role/powerTree/{id} [get]
func (r *_roleService) GetPowerTree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}

	var menus []*menu.SysMenu
	var powers []*uint
	err := dbctx.Use(c).Order("sort ASC").Find(&menus).Error
	if err != nil {
		global.GNA_LOG.Error("获取全部权限列表失败: " + err.Error())
		response.FailWithMessage("获取全部权限列表失败", c)
		return
	}

	powerTree := r.buildPowersTree(menus, 0)

	err = dbctx.Use(c).Table("sys_role_menu").
		Where("sys_role_id = ?", id).
		Pluck("sys_menu_id", &powers).Error
	if err != nil {
		global.GNA_LOG.Error("获取角色权限列表失败: " + err.Error())
		response.FailWithMessage("获取角色权限列表失败", c)
		return
	}

	response.OkWithData(PowerResponse{
		AllPowerTree: powerTree,
		RolePower:    powers,
	}, c)
}

func (r *_roleService) buildPowersTree(menus []*menu.SysMenu, parentID uint) []*AllPowerTree {
	var tree []*AllPowerTree
	for _, menu := range menus {
		if menu.ParentId == parentID {
			node := &AllPowerTree{
				Key:      menu.ID,
				Label:    menu.Title,
				Disabled: menu.Status == 0,
			}
			// 递归构建子菜单
			children := r.buildPowersTree(menus, menu.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// SetRolePower 设置角色权限（菜单ID列表）
// @Summary     分配角色菜单权限
// @Tags        角色
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SetPower true "角色 ID 与菜单 ID 列表"
// @Success     200 {object} response.Response
// @Router      /role/setRolePower [post]
func (r *_roleService) SetRolePower(c *gin.Context) {
	var setPowerReq SetPower
	err := c.ShouldBindJSON(&setPowerReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = dbctx.Use(c).Transaction(func(tx *gorm.DB) error {
		err := tx.Exec("DELETE FROM sys_role_menu WHERE sys_role_id = ?", setPowerReq.RoleId).Error
		if err != nil {
			return err
		}

		if len(setPowerReq.Powers) == 0 {
			return nil
		}
		var data []map[string]interface{}
		for _, powerId := range setPowerReq.Powers {
			data = append(data, map[string]interface{}{
				"sys_role_id": setPowerReq.RoleId,
				"sys_menu_id": powerId,
			})
		}
		return tx.Table("sys_role_menu").Create(&data).Error
	})

	if err != nil {
		global.GNA_LOG.Error("设置角色权限失败" + err.Error())
		response.FailWithMessage("设置角色权限失败", c)
		return
	}

	response.Ok(c)
}
