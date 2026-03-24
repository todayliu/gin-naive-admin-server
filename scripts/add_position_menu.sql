-- 职务管理菜单 SQL（组件路径：/system/position/index.vue）
-- 使用前请确认 sys_menu 中已存在「系统管理」目录（path='/system' 且 type='0'）
-- 若此前已插入 job-level 菜单，可执行下方「迁移」更新 component/path/name

INSERT INTO sys_menu (
  parent_id, type, path, name, component,
  title, title_i18n_key, icon, sort, status,
  create_time, update_time
)
SELECT
  id AS parent_id,
  '1' AS type,
  'position' AS path,
  'Position' AS name,
  '/system/position/index.vue' AS component,
  '职务管理' AS title,
  'routes.positionManagement' AS title_i18n_key,
  'ant-design:idcard-outlined' AS icon,
  36 AS sort,
  1 AS status,
  NOW() AS create_time,
  NOW() AS update_time
FROM sys_menu
WHERE path = '/system' AND type = '0' AND deleted_at IS NULL
LIMIT 1;

-- 迁移：从旧 job-level 菜单改为 position（按需执行）
-- UPDATE sys_menu SET path = 'position', name = 'Position', component = '/system/position/index.vue', title_i18n_key = 'routes.positionManagement', update_time = NOW()
-- WHERE path = 'job-level' AND deleted_at IS NULL;
