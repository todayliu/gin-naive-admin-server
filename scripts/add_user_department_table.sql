-- 用户-部门多对多关联表
-- 用于支持用户关联多个部门

CREATE TABLE IF NOT EXISTS sys_user_department (
  sys_user_id BIGINT UNSIGNED NOT NULL,
  sys_department_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (sys_user_id, sys_department_id),
  KEY idx_user_id (sys_user_id),
  KEY idx_department_id (sys_department_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户部门关联表';
