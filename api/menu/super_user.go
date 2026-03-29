package menu

import (
	"gin-admin-server/global"
	"strings"
)

// isMenuSuperUser 与 permission.IsSuperUser 规则一致（admin + 配置 super-role-codes）。
// 单独放在 menu 包，避免 menu → permission → menu 循环依赖。
func isMenuSuperUser(userID uint) bool {
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
