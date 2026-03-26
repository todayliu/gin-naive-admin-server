-- 参数配置页由 views/monitor/sys-config 迁至 views/system/sys-config 后执行一次
UPDATE sys_menu
SET component = 'system/sys-config/index.vue'
WHERE name = 'SysConfig'
  AND component = 'monitor/sys-config/index.vue'
  AND delete_time IS NULL;
