package role

import (
	"gin-admin-server/api/menu"
	"gin-admin-server/api/user"
	"gin-admin-server/global"
	"gin-admin-server/model/page_info"
)

type SysRole struct {
	global.GNA_MODEL
	Code        string `gorm:"column:code;unique;not null;comment:角色编码" json:"code" binding:"required" message:"角色编码不能为空"`
	Name        string `gorm:"column:name;50;not null;comment:角色名称" json:"name" binding:"required" message:"角色名称不能为空"`
	Description string `gorm:"column:description;comment:描述" json:"description"`
	Status      *uint  `gorm:"column:status;default:1;comment:角色状态" json:"status"`

	// 关联
	Menus []*menu.SysMenu `gorm:"many2many:sys_role_menu;" json:"-"`
	Users []*user.User    `gorm:"many2many:sys_user_role;" json:"-"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

type RolePageRequest struct {
	page_info.PageInfo
	Name   string
	Code   string
	Status *uint
}

type PowerResponse struct {
	AllPowerTree []*AllPowerTree `json:"allPowerTree"`
	RolePower    []*uint         `json:"rolePower"`
}

type AllPowerTree struct {
	Key      uint            `json:"key"`
	Label    string          `json:"label"`
	Disabled bool            `json:"disabled"`
	Children []*AllPowerTree `json:"children"`
	Sort     uint            `json:"sort"`
}

type SetPower struct {
	RoleId uint
	Powers []uint
}
