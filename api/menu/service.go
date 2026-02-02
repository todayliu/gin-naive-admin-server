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
	userId := jwt_util.GetUserID(c)
	routes, err := ms.GetMenuList(userId)
	if err != nil {
		response.FailWithMessage("获取用户菜单失败", c)
		return
	}

	response.OkWithData(routes, c)
}

func (ms *_menuService) GetMenuList(userId uint) ([]MenuResponse, error) {
	var menus []*SysMenu

	// 用户 → 角色 → 菜单
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
		return nil, err
	}
	menuTree := ms.buildMenuTree(menus, 0)

	var routes []MenuResponse
	for _, menu := range menuTree {
		route := ms.convertMenuToRoute(menu)
		routes = append(routes, route)
	}

	return routes, err
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

func (ms *_menuService) convertMenuToRoute(menu *SysMenu) MenuResponse {
	route := MenuResponse{
		Type:      menu.Type,
		Path:      menu.Path,
		Name:      menu.Name,
		Component: menu.Component,
		Sort:      menu.Sort,
		Status:    menu.Status,
		Meta: MenuMeta{
			Title:                menu.Title,
			Icon:                 menu.Icon,
			TitleI18nKey:         menu.TitleI18nKey,
			FixedInTabs:          menu.FixedInTabs,
			HideInMenu:           menu.HideInMenu,
			KeepAlive:            menu.KeepAlive,
			Link:                 menu.Link,
			LinkMode:             menu.LinkMode,
			NestedRouteRenderEnd: menu.NestedRouteRenderEnd,
		},
	}

	if menu.Redirect != "" {
		route.Redirect = &Redirect{
			Name: menu.Redirect,
		}
	}

	if len(menu.Children) > 0 {
		for _, child := range menu.Children {
			childRoute := ms.convertMenuToRoute(child)
			route.Children = append(route.Children, childRoute)
		}
	}

	return route
}
