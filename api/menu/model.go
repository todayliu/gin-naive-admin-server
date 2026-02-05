package menu

import (
	"gin-admin-server/api/user"
	"gin-admin-server/global"
)

// Menu 用户菜单表
type SysMenu struct {
	global.GNA_MODEL
	ParentId  uint    `gorm:"column:parent_id;not null;default:0;comment:父菜单ID" json:"parentId"`
	Type      string  `gorm:"column:type;not null;default:1;comment:菜单类型（0:目录  1:菜单  2:按钮）" json:"type"`
	Path      string  `gorm:"column:path;comment:路由路径" json:"path"`
	Name      string  `gorm:"column:name;unique;common:路由名称" json:"name"`
	Component string  `gorm:"column:component;comment:路由组件路径" json:"component"`
	Redirect  string  `gorm:"column:redirect;comment:重定向" json:"redirect"`
	Perms     *string `gorm:"column:perms;comment:授权标识" json:"perms"`

	//Meta 字段
	Title                string `gorm:"column:title;not null;comment:菜单标题" json:"title"  binding:"required" message:"菜单标题不能为空"`
	TitleI18nKey         string `gorm:"column:title_i18n_key;comment:国际化key" json:"titleI18nKey"`
	Icon                 string `gorm:"column:icon;comment:菜单标题" json:"icon"`
	FixedInTabs          bool   `gorm:"column:fixed_in_tabs;comment:固定在标签页" json:"fixedInTabs"`
	HideInMenu           bool   `gorm:"column:hide_in_menu;comment:不在菜单中显示" json:"hideInMenu"`
	KeepAlive            bool   `gorm:"column:keep_alive;comment:缓存配置" json:"keepAlive"`
	Link                 string `gorm:"column:link;comment:外链" json:"link"`
	LinkMode             string `gorm:"column:link_mode;comment:链接模式" json:"linkMode"`
	NestedRouteRenderEnd bool   `gorm:"column:nested_route_render_end;comment:是否在当前路由层级结束渲染" json:"nestedRouteRenderEnd"`

	//系统字段
	Sort   uint `gorm:"column:sort;default:0;comment:排序" json:"sort"`
	Status uint `gorm:"column:status;default:1;comment:状态（1:启用  2:禁用）" json:"status"`
	// 关联
	Children []*SysMenu `gorm:"-" json:"children,omitempty"`
	Roles    []*SysRole `gorm:"many2many:sys_role_menu;" json:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}

// SysRole 用户角色表
type SysRole struct {
	global.GNA_MODEL
	Code        string `gorm:"column:code;unique;not null;comment:角色编码" json:"code"`
	Name        string `gorm:"column:name;50;not null;comment:角色名称" json:"name"`
	Description string `gorm:"column:description;comment:描述" json:"description"`
	DataScope   int8   `gorm:"column:data_scope;default:1;comment:数据权限" json:"dataScope"`
	Status      int8   `gorm:"column:status;default:1;comment:角色状态" json:"status"`

	// 关联
	Menus []*SysMenu   `gorm:"many2many:sys_role_menu;" json:"menus,omitempty"`
	Users []*user.User `gorm:"many2many:sys_user_role;" json:"users,omitempty"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

type MenuResponse struct {
	Type      string         `json:"type"`
	Path      string         `json:"path"`
	Name      string         `json:"name"`
	Component string         `json:"component"`
	Redirect  *Redirect      `json:"redirect"`
	Perms     *string        `json:"perms"`
	Sort      uint           `json:"sort"`
	Status    uint           `json:"status"`
	Meta      MenuMeta       `json:"meta"`
	Children  []MenuResponse `json:"children"`
}

type MenuMeta struct {
	Title                string `json:"title"`
	Icon                 string `json:"icon"`
	TitleI18nKey         string `json:"titleI18nKey"`
	FixedInTabs          bool   `json:"fixedInTabs"`
	HideInMenu           bool   `json:"hideInMenu"`
	KeepAlive            bool   `json:"keepAlive"`
	Link                 string `json:"link"`
	LinkMode             string `json:"linkMode"`
	NestedRouteRenderEnd bool   `json:"nestedRouteRenderEnd"`
}
type Redirect struct {
	Name string `json:"name"`
}
