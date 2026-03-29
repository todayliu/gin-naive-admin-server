// Package role 角色模块 API
package role

import (
	"gin-admin-server/api/menu"
	"gin-admin-server/api/user"
	"gin-admin-server/global"
	"gin-admin-server/model/page_info"
)

// SysRole 角色表
type SysRole struct {
	global.GNA_MODEL
	Code        string `gorm:"column:code;unique;not null;comment:角色编码" json:"code" binding:"required" message:"角色编码不能为空"`
	Name        string `gorm:"column:name;50;not null;comment:角色名称" json:"name" binding:"required" message:"角色名称不能为空"`
	Description string `gorm:"column:description;comment:描述" json:"description"`
	Status      *uint  `gorm:"column:status;default:1;comment:角色状态" json:"status"`

	// 关联
	Menus []*menu.SysMenu `gorm:"many2many:sys_role_menu;" json:"-"`
	Users []*user.SysUser `gorm:"many2many:sys_user_role;" json:"-"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

// RoleListFilters 角色列表/导出共用的查询条件（不含分页；导出接口只绑定此结构）
type RoleListFilters struct {
	Name   string `form:"name"`
	Code   string `form:"code"`
	Status *uint  `form:"status"`
}

// RolePageRequest 角色分页查询请求
type RolePageRequest struct {
	page_info.PageInfo
	RoleListFilters
}

// PowerResponse 角色权限树响应
type PowerResponse struct {
	AllPowerTree []*AllPowerTree `json:"allPowerTree"`
	RolePower    []*uint         `json:"rolePower"`
}

// AllPowerTree 权限树节点
type AllPowerTree struct {
	Key      uint            `json:"key"`
	Label    string          `json:"label"`
	Disabled bool            `json:"disabled"`
	Children []*AllPowerTree `json:"children"`
	Sort     uint            `json:"sort"`
}

// SetPower 设置角色权限请求
type SetPower struct {
	RoleId uint   `json:"roleId" binding:"required"`
	Powers []uint `json:"powers"`
}
