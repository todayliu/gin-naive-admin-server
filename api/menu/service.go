package menu

import (
	"fmt"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/jwt_util"
	"gin-admin-server/utils/validator"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type _menuService struct{}

var MenuService = new(_menuService)

// InitMenuList 根据当前用户获取可访问的菜单路由（用于动态路由）
// @Summary     当前用户菜单路由树
// @Description 用于前端动态路由注册；需登录。
// @Tags        菜单
// @Produce     json
// @Security    AccessToken
// @Success     200 {object} response.Response
// @Router      /menu/router [get]
func (ms *_menuService) InitMenuList(c *gin.Context) {
	userId := jwt_util.GetUserID(c)
	routes, err := ms.GetMenuList(userId)
	if err != nil {
		response.FailWithMessage("获取用户菜单失败", c)
		return
	}

	response.OkWithData(routes, c)
}

// GetMenuList 根据用户ID获取其可访问的菜单树（用户→角色→菜单）
// 超级管理员：返回库中全部启用目录/页面（type 0、1），与「全部接口权限码」一致，不要求 sys_role_menu 逐条授权。
func (ms *_menuService) GetMenuList(userId uint) ([]MenuResponse, error) {
	var menus []*SysMenu
	var err error

	if IsSuperUser(userId) {
		err = global.GNA_DB.Where("status = ? AND type IN ?", 1, []string{"0", "1"}).Order("sort ASC").Find(&menus).Error
	} else {
		// 用户 → 角色 → 菜单
		err = global.GNA_DB.Unscoped().
			Table("sys_user_role ur").
			Select("DISTINCT m.*").
			Joins("JOIN sys_role_menu rm ON ur.sys_role_id = rm.sys_role_id").
			Joins("JOIN sys_menu m ON rm.sys_menu_id = m.id").
			Where("ur.sys_user_id = ? AND m.status = 1 AND m.type IN (0, 1)", userId).
			Order("m.sort ASC").
			Find(&menus).Error
	}

	if err != nil {
		global.GNA_LOG.Error("获取用户菜单失败: " + err.Error())
		return nil, err
	}
	menuTree := ms.BuildMenuTree(menus, 0)

	var routes []MenuResponse
	for _, menu := range menuTree {
		route := ms.convertMenuToRoute(menu)
		routes = append(routes, route)
	}

	return routes, err
}

// BuildMenuTree 递归构建菜单树
func (ms *_menuService) BuildMenuTree(menus []*SysMenu, parentID uint) []*SysMenu {
	var tree []*SysMenu
	for _, menu := range menus {
		if menu.ParentId == parentID {
			// 递归构建子菜单
			children := ms.BuildMenuTree(menus, menu.ID)
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
		Perms:     menu.Perms,
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
	} else {
		route.Redirect = nil
	}

	if len(menu.Children) > 0 {
		for _, child := range menu.Children {
			childRoute := ms.convertMenuToRoute(child)
			route.Children = append(route.Children, childRoute)
		}
	}

	return route
}

// GetAllMenuList 获取全部菜单树（用于菜单管理）
// @Summary     全部菜单树
// @Tags        菜单
// @Produce     json
// @Security    AccessToken
// @Success     200 {object} response.Response
// @Router      /menu/list [get]
func (ms *_menuService) GetAllMenuList(c *gin.Context) {
	var menus []*SysMenu
	err := global.GNA_DB.Order("sort ASC").Find(&menus).Error
	if err != nil {
		global.GNA_LOG.Error("获取菜单列表失败: " + err.Error())
		response.FailWithMessage("获取菜单列表失败", c)
		return
	}

	menuTree := ms.BuildMenuTree(menus, 0)

	response.OkWithData(menuTree, c)
}

// fillMenuNameIfEmpty 路由 name 唯一；未填时由 perms 或标题生成，避免 Create 失败
func (ms *_menuService) fillMenuNameIfEmpty(menu *SysMenu) {
	if strings.TrimSpace(menu.Name) != "" {
		return
	}
	var base string
	if menu.Perms != nil && strings.TrimSpace(*menu.Perms) != "" {
		base = strings.ReplaceAll(strings.TrimSpace(*menu.Perms), ":", "_")
		base = strings.ReplaceAll(base, "/", "_")
	} else if strings.TrimSpace(menu.Title) != "" {
		base = strings.ReplaceAll(strings.TrimSpace(menu.Title), " ", "_")
	} else {
		base = "menu"
	}
	if len(base) > 100 {
		base = base[:100]
	}
	for i := 0; i < 200; i++ {
		candidate := base
		if i > 0 {
			candidate = fmt.Sprintf("%s_%d", base, i)
		}
		var cnt int64
		global.GNA_DB.Model(&SysMenu{}).Where("name = ?", candidate).Count(&cnt)
		if cnt == 0 {
			menu.Name = candidate
			return
		}
	}
	menu.Name = fmt.Sprintf("%s_%d", base, time.Now().UnixMilli())
}

// AddMenu 新增菜单（不含 id）
// @Summary     新增菜单
// @Tags        菜单
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysMenu true "菜单实体（勿传 id）"
// @Success     200 {object} response.Response
// @Router      /menu/add [post]
func (ms *_menuService) AddMenu(c *gin.Context) {
	var menu SysMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, menu), c)
		return
	}
	menu.ID = 0
	ms.fillMenuNameIfEmpty(&menu)
	if err := global.GNA_DB.Create(&menu).Error; err != nil {
		global.GNA_LOG.Error("新增菜单失败", zap.Error(err))
		response.FailWithMessage("新增菜单失败", c)
		return
	}
	response.Ok(c)
}

// UpdateMenu 修改菜单（必须带 id）
// @Summary     编辑菜单
// @Tags        菜单
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysMenu true "菜单实体（含 id）"
// @Success     200 {object} response.Response
// @Router      /menu/edit [put]
func (ms *_menuService) UpdateMenu(c *gin.Context) {
	var menu SysMenu
	err := c.ShouldBindJSON(&menu)

	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, menu)
		global.GNA_LOG.Error(errMessage)
		response.FailWithMessage(errMessage, c)
		return
	}
	if menu.ID == 0 {
		response.FailWithMessage("缺少菜单 id，请使用 POST /api/menu/add 新增", c)
		return
	}
	err = global.GNA_DB.Save(&menu).Error

	if err != nil {
		global.GNA_LOG.Error("菜单修改失败：" + err.Error())
		response.FailWithMessage("菜单修改失败", c)
		return
	}

	response.Ok(c)
}

// DeleteMenu 删除菜单（永久删除，同时删除角色菜单关联）
// @Summary     删除菜单
// @Tags        菜单
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "菜单 ID"
// @Success     200 {object} response.Response
// @Router      /menu/delete/{id} [delete]
func (ms *_menuService) DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}

	err := global.GNA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM sys_role_menu WHERE sys_menu_id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Where("id = ?", id).Delete(&SysMenu{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.GNA_LOG.Error("删除菜单失败：" + err.Error())
		response.FailWithMessage("删除菜单失败", c)
		return
	}

	response.OkWithMessage("删除菜单成功", c)
}
