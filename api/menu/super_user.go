package menu

import (
	"gin-admin-server/global"
	"strings"
)

// IsSuperUser 是否拥有超级角色（admin 或配置中的 super-role-codes）；用于菜单、数据范围、中间件等。
func IsSuperUser(userID uint) bool {
	codes := resolvedMenuSuperRoleCodes()
	var n int64
	err := global.GNA_DB.Table("sys_user_role ur").
		Joins("JOIN sys_role r ON r.id = ur.sys_role_id AND r.delete_time IS NULL").
		Where("ur.sys_user_id = ? AND r.code IN ?", userID, codes).
		Limit(1).
		Count(&n).Error
	if err != nil {
		return false
	}
	return n > 0
}

func resolvedMenuSuperRoleCodes() []string {
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
