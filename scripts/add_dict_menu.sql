-- 字典管理菜单 SQL
-- 使用前请确认 sys_menu 中已存在「系统管理」目录（path='/system' 或 name='System'）
-- 将字典菜单插入到系统管理目录下

INSERT INTO sys_menu (
  parent_id, type, path, name, component,
  title, title_i18n_key, icon, sort, status,
  create_time, update_time
)
SELECT
  id AS parent_id,
  '1' AS type,
  'dict' AS path,
  'Dict' AS name,
  '/system/dict/index.vue' AS component,
  '字典管理' AS title,
  'routes.dictManagement' AS title_i18n_key,
  'ant-design:book-outlined' AS icon,
  40 AS sort,
  1 AS status,
  NOW() AS create_time,
  NOW() AS update_time
FROM sys_menu
WHERE path = '/system' AND type = '0' AND deleted_at IS NULL
LIMIT 1;

-- 若需为超级管理员角色授权，请先执行上述插入，获取新菜单 id，再执行：
-- INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) 
-- SELECT id, (SELECT id FROM sys_menu WHERE name='Dict' AND deleted_at IS NULL ORDER BY id DESC LIMIT 1) 
-- FROM sys_role WHERE name = '超级管理员' LIMIT 1;
