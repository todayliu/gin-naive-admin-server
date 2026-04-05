package menu

import (
	"gin-admin-server/global"
)

// LoadUserButtonPermCodes 从菜单表 sys_menu.perms 汇总权限标识：普通用户仅含其角色在 sys_role_menu 中关联到的菜单；超管为库中全部非空 perms（与 InitMenuList 超管逻辑一致，数据来源为菜单表而非写死列表）。
func LoadUserButtonPermCodes(userID uint) []string {
	if userID == 0 {
		return nil
	}
	if IsSuperUser(userID) {
		return loadAllDistinctPermsFromDB()
	}
	return loadUserPermCodesFromRoleMenus(userID)
}

func loadUserPermCodesFromRoleMenus(userID uint) []string {
	var perms []string
	err := global.GNA_DB.Table("sys_user_role ur").
		Select("DISTINCT m.perms").
		Joins("JOIN sys_role_menu rm ON rm.sys_role_id = ur.sys_role_id").
		Joins("JOIN sys_menu m ON m.id = rm.sys_menu_id AND m.delete_time IS NULL").
		Where("ur.sys_user_id = ? AND m.perms IS NOT NULL AND m.perms != ?", userID, "").
		Pluck("m.perms", &perms).Error
	if err != nil {
		global.GNA_LOG.Sugar().Errorf("loadUserPermCodesFromRoleMenus: %v", err)
		return nil
	}
	return perms
}

func loadAllDistinctPermsFromDB() []string {
	var perms []string
	err := global.GNA_DB.Model(&SysMenu{}).
		Select("DISTINCT perms").
		Where("perms IS NOT NULL AND perms != ?", "").
		Pluck("perms", &perms).Error
	if err != nil {
		global.GNA_LOG.Sugar().Errorf("loadAllDistinctPermsFromDB: %v", err)
		return nil
	}
	return perms
}
