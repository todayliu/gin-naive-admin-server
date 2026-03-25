package permission

import (
	"gin-admin-server/api/menu"
	"gin-admin-server/global"
	"strings"

	"gorm.io/gorm"
)

// resolvedSuperRoleCodes 用于判断是否超级用户：角色编码 admin 恒在其中，并与配置 super-role-codes 合并去重
func resolvedSuperRoleCodes() []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, 8)
	add := func(c string) {
		c = strings.TrimSpace(c)
		if c == "" {
			return
		}
		if _, ok := seen[c]; ok {
			return
		}
		seen[c] = struct{}{}
		out = append(out, c)
	}
	add("admin")
	for _, c := range global.GNA_CONFIG.Security.SuperRoleCodes {
		add(c)
	}
	return out
}

// IsSuperUser 是否拥有超级角色（admin 或配置中的 super-role-codes）
func IsSuperUser(userID uint) bool {
	codes := resolvedSuperRoleCodes()
	var n int64
	err := global.GNA_DB.Table("sys_user_role ur").
		Joins("JOIN sys_role r ON r.id = ur.sys_role_id AND r.delete_time IS NULL"). // 与 global.GNA_MODEL 列名一致
		Where("ur.sys_user_id = ? AND r.code IN ?", userID, codes).
		Limit(1).
		Count(&n).Error
	if err != nil {
		global.GNA_LOG.Sugar().Errorf("IsSuperUser: %v", err)
		return false
	}
	return n > 0
}

// LoadUserRoleCodes 用户拥有的角色 code 列表（如 admin）
func LoadUserRoleCodes(userID uint) ([]string, error) {
	var codes []string
	err := global.GNA_DB.Table("sys_user_role ur").
		Joins("JOIN sys_role r ON r.id = ur.sys_role_id AND r.delete_time IS NULL").
		Where("ur.sys_user_id = ?", userID).
		Pluck("r.code", &codes).Error
	return codes, err
}

// LoadUserEffectivePermCodes 登录/会话用：超管返回全部接口码，否则返回菜单按钮 perms
func LoadUserEffectivePermCodes(userID uint) []string {
	if IsSuperUser(userID) {
		return AllRegisteredCodes()
	}
	codes, err := LoadUserPermCodes(userID)
	if err != nil {
		return nil
	}
	return codes
}

// LoadUserPermCodes 从菜单按钮 perms 字段加载用户权限标识
func LoadUserPermCodes(userID uint) ([]string, error) {
	var perms []string
	err := global.GNA_DB.Table("sys_user_role ur").
		Select("DISTINCT m.perms").
		Joins("JOIN sys_role_menu rm ON rm.sys_role_id = ur.sys_role_id").
		Joins("JOIN sys_menu m ON m.id = rm.sys_menu_id AND m.delete_time IS NULL").
		Where("ur.sys_user_id = ? AND m.perms IS NOT NULL AND m.perms != ?", userID, "").
		Pluck("m.perms", &perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}

// UserHasPermCode 检查用户是否拥有某权限码（含角色菜单 perms）
func UserHasPermCode(userID uint, code string) bool {
	if code == "" {
		return true
	}
	if IsSuperUser(userID) {
		return true
	}
	codes, err := LoadUserPermCodes(userID)
	if err != nil {
		return false
	}
	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}

// SeedMenuButtonPermsIfNeeded 若库中尚无任何带 perms 的菜单，则按 permission 注册表为 admin 角色生成按钮权限
func SeedMenuButtonPermsIfNeeded(db *gorm.DB) {
	if db == nil {
		return
	}
	var cnt int64
	db.Model(&menu.SysMenu{}).Where("perms IS NOT NULL AND perms != ?", "").Count(&cnt)
	if cnt > 0 {
		return
	}
	var adminID uint
	if err := db.Table("sys_role").Select("id").Where("code = ? AND delete_time IS NULL", "admin").Scan(&adminID).Error; err != nil || adminID == 0 {
		return
	}
	codes := AllRegisteredCodes()
	for i, code := range codes {
		name := "apiPerm_" + strings.ReplaceAll(code, ":", "_")
		title := code
		p := code
		parentID := resolveButtonParentID(db, code)
		if parentID == 0 {
			parentID = findSystemDirectoryID(db)
		}
		m := menu.SysMenu{
			ParentId:   parentID,
			Type:       "2",
			Path:       "",
			Name:       name,
			Title:      title,
			Sort:       uint(i + 1),
			Status:     1,
			HideInMenu: true,
			Perms:      &p,
		}
		if err := db.Create(&m).Error; err != nil {
			continue
		}
		_ = db.Exec("INSERT IGNORE INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES (?, ?)", adminID, m.ID).Error
	}
}
