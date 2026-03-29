package permission

import "gin-admin-server/global"

// AllRegisteredCodes 所有受控接口对应的权限标识（用于超管登录时下发的 codes）
func AllRegisteredCodes() []string {
	m := buildRoutePermMap(global.GNA_CONFIG.Router.RouterPrefix)
	seen := make(map[string]struct{}, len(m))
	out := make([]string, 0, len(m))
	for _, code := range m {
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		out = append(out, code)
	}
	return out
}

// RequiredCode 根据路由模板返回所需权限码；第二个返回 false 表示未配置强制权限（仅 JWT 即可）
func RequiredCode(method, fullPath string) (string, bool) {
	m := buildRoutePermMap(global.GNA_CONFIG.Router.RouterPrefix)
	code, ok := m[method+" "+fullPath]
	return code, ok
}

func buildRoutePermMap(prefix string) map[string]string {
	p := prefix
	if p == "" {
		p = "/api"
	}
	m := make(map[string]string)
	put := func(method, path string, code string) {
		m[method+" "+p+path] = code
	}

	// 部门（列表 GET 不校验权限码，仅 JWT；下同）
	put("GET", "/department/export", "system:department:export")
	put("PUT", "/department/edit", "system:department:edit")
	put("DELETE", "/department/delete/:id", "system:department:delete")

	// 用户
	put("GET", "/user/query/:id", "system:user:query")
	put("GET", "/user/roles/:id", "system:user:roles")
	put("GET", "/user/export", "system:user:export")
	put("GET", "/user/import-template", "system:user:import")
	put("POST", "/user/import", "system:user:import")
	put("POST", "/user/add", "system:user:add")
	put("PUT", "/user/edit", "system:user:edit")
	put("DELETE", "/user/delete/:id", "system:user:delete")

	// 角色
	put("GET", "/role/export", "system:role:export")
	put("GET", "/role/import-template", "system:role:import")
	put("POST", "/role/import", "system:role:import")
	put("POST", "/role/add", "system:role:add")
	put("GET", "/role/query/:id", "system:role:query")
	put("PUT", "/role/edit", "system:role:edit")
	put("DELETE", "/role/delete/:id", "system:role:delete")
	put("GET", "/role/powerTree/:id", "system:role:power")
	put("POST", "/role/setRolePower", "system:role:power")

	// 菜单（router 由中间件白名单放行，此处仅管理端菜单维护）
	put("POST", "/menu/add", "system:menu:add")
	put("PUT", "/menu/edit", "system:menu:edit")
	put("DELETE", "/menu/delete/:id", "system:menu:delete")

	// 字典
	put("POST", "/dict/type/add", "system:dict:type:add")
	put("PUT", "/dict/type/edit", "system:dict:type:edit")
	put("DELETE", "/dict/type/delete/:id", "system:dict:type:delete")
	put("POST", "/dict/data/add", "system:dict:data:add")
	put("PUT", "/dict/data/edit", "system:dict:data:edit")
	put("DELETE", "/dict/data/delete/:id", "system:dict:data:delete")

	// 职务
	put("GET", "/position/export", "system:position:export")
	put("GET", "/position/import-template", "system:position:import")
	put("POST", "/position/import", "system:position:import")
	put("GET", "/position/query/:id", "system:position:query")
	put("POST", "/position/add", "system:position:add")
	put("PUT", "/position/edit", "system:position:edit")
	put("DELETE", "/position/delete/:id", "system:position:delete")

	// 个人中心、字典按编码查询、菜单路由：中间件白名单（仅需登录）

	// 文件
	put("POST", "/file/upload", "system:file:upload")

	// 日志（列表 GET 不校验权限码）

	// 参数配置
	put("GET", "/config/export", "system:config:export")
	put("GET", "/config/import-template", "system:config:import")
	put("POST", "/config/import", "system:config:import")
	put("POST", "/config/add", "system:config:add")
	put("PUT", "/config/edit", "system:config:edit")

	// 仪表盘
	put("GET", "/dashboard/stats", "system:dashboard:stats")

	return m
}
