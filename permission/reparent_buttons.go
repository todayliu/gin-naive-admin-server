package permission

import (
	"gin-admin-server/api/menu"
	"gin-admin-server/global"
	"strings"

	"gorm.io/gorm"
)

// parentPathForPermCode 权限码对应的「页面菜单」path（type=1）；空字符串表示挂在系统根目录下
func parentPathForPermCode(code string) string {
	switch {
	case strings.HasPrefix(code, "system:department:"):
		return "department"
	case strings.HasPrefix(code, "system:user:"):
		return "user"
	case strings.HasPrefix(code, "system:role:"):
		return "role"
	case strings.HasPrefix(code, "system:menu:"):
		return "menu"
	case strings.HasPrefix(code, "system:dict:"):
		return "dict"
	case strings.HasPrefix(code, "system:position:"):
		return "position"
	case strings.HasPrefix(code, "system:file:"):
		return ""
	case strings.HasPrefix(code, "system:log:login:"):
		return "log-login"
	case strings.HasPrefix(code, "system:log:oper:"):
		return "log-oper"
	case strings.HasPrefix(code, "system:config:"):
		return "sys-config"
	case strings.HasPrefix(code, "system:dashboard:"):
		return "home"
	default:
		return ""
	}
}

func findSystemDirectoryID(db *gorm.DB) uint {
	if db == nil {
		return 0
	}
	var id uint
	for _, p := range []string{"/system", "system"} {
		db.Table("sys_menu").Select("id").Where("type = ? AND path = ? AND delete_time IS NULL", "0", p).Order("id ASC").Limit(1).Scan(&id)
		if id > 0 {
			return id
		}
	}
	return 0
}

// findLeafMenuID 按页面菜单 path 查找 type=1 的菜单 id；可传多个候选 path（如 position / job-level）
func findLeafMenuID(db *gorm.DB, paths ...string) uint {
	if db == nil {
		return 0
	}
	for _, p := range paths {
		if p == "" {
			continue
		}
		var id uint
		db.Table("sys_menu").Select("id").Where("type = ? AND path = ? AND delete_time IS NULL", "1", p).Order("id ASC").Limit(1).Scan(&id)
		if id > 0 {
			return id
		}
	}
	return 0
}

// resolveButtonParentID 计算接口权限按钮应挂载的父菜单 id
func resolveButtonParentID(db *gorm.DB, permCode string) uint {
	if db == nil {
		return 0
	}
	p := parentPathForPermCode(permCode)
	sysRoot := findSystemDirectoryID(db)

	switch p {
	case "position":
		if id := findLeafMenuID(db, "position", "job-level"); id > 0 {
			return id
		}
	case "home":
		if id := findLeafMenuID(db, "home"); id > 0 {
			return id
		}
	case "log-login":
		if id := findLeafMenuID(db, "log-login"); id > 0 {
			return id
		}
	case "log-oper":
		if id := findLeafMenuID(db, "log-oper"); id > 0 {
			return id
		}
	case "sys-config":
		if id := findLeafMenuID(db, "sys-config", "config"); id > 0 {
			return id
		}
	case "":
		if sysRoot > 0 {
			return sysRoot
		}
		return 0
	default:
		if id := findLeafMenuID(db, p); id > 0 {
			return id
		}
	}

	if sysRoot > 0 {
		return sysRoot
	}
	return 0
}

// ReparentAPIPermissionButtons 将 apiPerm_* 接口权限按钮挂到对应业务菜单下；可重复执行（用于修正历史数据）
func ReparentAPIPermissionButtons(db *gorm.DB) {
	if db == nil {
		return
	}
	for _, code := range AllRegisteredCodes() {
		name := "apiPerm_" + strings.ReplaceAll(code, ":", "_")
		pid := resolveButtonParentID(db, code)
		if pid == 0 {
			continue
		}
		if err := db.Model(&menu.SysMenu{}).
			Where("name = ? AND type = ? AND delete_time IS NULL", name, "2").
			Update("parent_id", pid).Error; err != nil {
			global.GNA_LOG.Sugar().Warnf("reparent button %s: %v", name, err)
		}
	}
}
