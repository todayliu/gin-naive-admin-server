-- 可选：在「系统管理」下增加登录日志、操作日志、参数配置菜单
-- 若 @parent 为空，请先确认存在 path='system' 且 type='0' 的目录行

SET @parent := (SELECT id FROM sys_menu WHERE path = 'system' AND type = '0' AND delete_time IS NULL ORDER BY id LIMIT 1);

INSERT INTO sys_menu (
  parent_id, type, path, name, component, title, title_i18n_key, icon, sort, status,
  hide_in_menu, fixed_in_tabs, keep_alive, nested_route_render_end
)
SELECT @parent, '1', 'log-login', 'LogLogin', 'monitor/log-login/index.vue', '登录日志', 'routes.loginLog', 'mdi:text-box-search-outline', 80, 1,
  0, 0, 1, 0
WHERE @parent IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM sys_menu WHERE name = 'LogLogin' AND delete_time IS NULL);

INSERT INTO sys_menu (
  parent_id, type, path, name, component, title, title_i18n_key, icon, sort, status,
  hide_in_menu, fixed_in_tabs, keep_alive, nested_route_render_end
)
SELECT @parent, '1', 'log-oper', 'LogOper', 'monitor/log-oper/index.vue', '操作日志', 'routes.operLog', 'mdi:clipboard-text-outline', 81, 1,
  0, 0, 1, 0
WHERE @parent IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM sys_menu WHERE name = 'LogOper' AND delete_time IS NULL);

INSERT INTO sys_menu (
  parent_id, type, path, name, component, title, title_i18n_key, icon, sort, status,
  hide_in_menu, fixed_in_tabs, keep_alive, nested_route_render_end
)
SELECT @parent, '1', 'sys-config', 'SysConfig', 'system/sys-config/index.vue', '参数配置', 'routes.sysConfig', 'mdi:cog-outline', 82, 1,
  0, 0, 1, 0
WHERE @parent IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM sys_menu WHERE name = 'SysConfig' AND delete_time IS NULL);

INSERT IGNORE INTO sys_role_menu (sys_role_id, sys_menu_id)
SELECT r.id, m.id
FROM sys_role r
JOIN sys_menu m ON m.name IN ('LogLogin', 'LogOper', 'SysConfig') AND m.delete_time IS NULL
WHERE r.code = 'admin' AND r.delete_time IS NULL;

-- 若此前已插入旧路径，可执行：
-- UPDATE sys_menu SET component = 'monitor/log-login/index.vue' WHERE name = 'LogLogin' AND component LIKE 'system/log-login%' AND delete_time IS NULL;
-- UPDATE sys_menu SET component = 'monitor/log-oper/index.vue' WHERE name = 'LogOper' AND component LIKE 'system/log-oper%' AND delete_time IS NULL;
-- UPDATE sys_menu SET component = 'system/sys-config/index.vue' WHERE name = 'SysConfig' AND component = 'monitor/sys-config/index.vue' AND delete_time IS NULL;
