package menu

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/jwt_util"
	"sort"

	"github.com/gin-gonic/gin"
)

type _menuService struct{}

var MenuService = new(_menuService)

func (ms *_menuService) InitMenuList(c *gin.Context) {
	//var menus []*SysMenu
	ms.getSysRole(c)
	//response.OkWithMessage("获取菜单成功", c)
}

func (ms *_menuService) getSysRole(c *gin.Context) {
	userId := jwt_util.GetUserID(c)

	//var userRoleIDs []uint
	//err := global.GNA_DB.Table("sys_user_role").Select("sys_role_id").Where("user_id = ?", userId).Pluck("sys_role_id", &userRoleIDs).Error
	//if err != nil {
	//	global.GNA_LOG.Error("获取用户角色信息失败: " + err.Error())
	//	response.FailWithMessage("获取用户角色信息失败", c)
	//	return
	//}
	//if len(userRoleIDs) == 0 {
	//	response.OkWithData([]*SysMenu{}, c)
	//	return
	//}
	//
	//var menuIDs []uint
	//err = global.GNA_DB.Table("sys_role_menu").Select("DISTINCT sys_menu_id").Where("sys_role_id IN ?", userRoleIDs).Pluck("sys_menu_id", &menuIDs).Error
	//if err != nil {
	//	global.GNA_LOG.Error("获取用户角色菜单失败: " + err.Error())
	//	response.FailWithMessage("获取用户角色菜单失败", c)
	//	return
	//}
	//if len(menuIDs) == 0 {
	//	response.OkWithData([]*SysMenu{}, c)
	//	return
	//}
	//
	//var menus []*SysMenu
	//err = global.GNA_DB.Where("id IN ? AND status = 1", menuIDs).Order("sort ASC").Find(&menus).Error
	//if err != nil {
	//	global.GNA_LOG.Error("获取用户菜单失败: " + err.Error())
	//	response.FailWithMessage("获取用户菜单失败", c)
	//	return
	//}
	//if len(menus) == 0 {
	//	response.OkWithData([]*SysMenu{}, c)
	//	return
	//}

	var menus []*SysMenu

	// 一条SQL完成：用户 → 角色 → 菜单
	err := global.GNA_DB.Unscoped().
		Table("sys_user_role ur").
		Select("DISTINCT m.*").
		Joins("JOIN sys_role_menu rm ON ur.sys_role_id = rm.sys_role_id").
		Joins("JOIN sys_menu m ON rm.sys_menu_id = m.id").
		Where("ur.user_id = ? AND m.status = 1", userId).
		Order("m.sort ASC").
		Find(&menus).Error

	if err != nil {
		global.GNA_LOG.Error("获取用户菜单失败: " + err.Error())
		response.FailWithMessage("获取用户菜单失败", c)
		return
	}
	menuTree := ms.buildMenuTree(menus, 0)
	response.OkWithData(menuTree, c)
}

func (ms *_menuService) buildMenuTree(menus []*SysMenu, parentID uint) []*SysMenu {
	var tree []*SysMenu
	for _, menu := range menus {
		if menu.ParentId == parentID {
			// 递归构建子菜单
			children := ms.buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menu.Children = children
			}
			tree = append(tree, menu)
		}
	}
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].Sort < tree[j].Sort
	})

	return tree
}
