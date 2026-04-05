package devform

import (
	"fmt"
	"gin-admin-server/global"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var (
	validName   = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)
	reservedCol = map[string]struct{}{
		"id": {}, "create_by": {}, "update_by": {}, "delete_by": {},
		"create_time": {}, "update_time": {}, "delete_time": {},
	}
	allowedDbTypes = map[string]struct{}{
		"varchar": {}, "text": {}, "int": {}, "bigint": {}, "tinyint": {},
		"datetime": {}, "date": {}, "decimal": {},
	}
)

func validateIdent(name string) bool { return validName.MatchString(name) }

func validateDBType(t string) bool {
	_, ok := allowedDbTypes[strings.ToLower(t)]
	return ok
}

func escapeSQLComment(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "'", "''"), "\n", " ")
}

func columnDDL(f SysOnlineFormField) (string, error) {
	dt := strings.ToLower(f.DbType)
	null := "NOT NULL"
	if f.Nullable {
		null = "NULL DEFAULT NULL"
	}
	cm := escapeSQLComment(f.Comment)
	comment := ""
	if cm != "" {
		comment = fmt.Sprintf(" COMMENT '%s'", cm)
	}
	switch dt {
	case "varchar":
		L := f.Length
		if L <= 0 || L > 2000 {
			L = 255
		}
		return fmt.Sprintf("`%s` VARCHAR(%d) %s%s", f.ColumnName, L, null, comment), nil
	case "text":
		return fmt.Sprintf("`%s` TEXT %s%s", f.ColumnName, null, comment), nil
	case "int":
		return fmt.Sprintf("`%s` INT %s%s", f.ColumnName, null, comment), nil
	case "bigint":
		return fmt.Sprintf("`%s` BIGINT %s%s", f.ColumnName, null, comment), nil
	case "tinyint":
		return fmt.Sprintf("`%s` TINYINT %s%s", f.ColumnName, null, comment), nil
	case "datetime":
		return fmt.Sprintf("`%s` DATETIME(3) %s%s", f.ColumnName, null, comment), nil
	case "date":
		return fmt.Sprintf("`%s` DATE %s%s", f.ColumnName, null, comment), nil
	case "decimal":
		sc := f.DecimalScale
		if sc < 0 || sc > 8 {
			sc = 2
		}
		return fmt.Sprintf("`%s` DECIMAL(18,%d) %s%s", f.ColumnName, sc, null, comment), nil
	default:
		return "", fmt.Errorf("不支持的 db_type: %s", f.DbType)
	}
}

func buildCreateTableSQL(table string, fields []SysOnlineFormField) (string, error) {
	var b strings.Builder
	b.WriteString("CREATE TABLE IF NOT EXISTS `")
	b.WriteString(table)
	b.WriteString("` (\n")
	b.WriteString("  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,\n")
	b.WriteString("  `create_by` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人用户ID',\n")
	b.WriteString("  `update_by` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新人用户ID',\n")
	b.WriteString("  `delete_by` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '删除人用户ID',\n")
	b.WriteString("  `create_time` DATETIME(3) NULL COMMENT '创建时间',\n")
	b.WriteString("  `update_time` DATETIME(3) NULL COMMENT '更新时间',\n")
	b.WriteString("  `delete_time` DATETIME(3) NULL COMMENT '软删除时间',\n")
	for _, f := range fields {
		if _, skip := reservedCol[strings.ToLower(f.ColumnName)]; skip {
			continue
		}
		col, err := columnDDL(f)
		if err != nil {
			return "", err
		}
		b.WriteString("  ")
		b.WriteString(col)
		b.WriteString(",\n")
	}
	b.WriteString("  PRIMARY KEY (`id`),\n")
	b.WriteString("  KEY `idx_")
	b.WriteString(table)
	b.WriteString("_delete_time` (`delete_time`)\n")
	b.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	return b.String(), nil
}

func tableExists(db *gorm.DB, schema, table string) (bool, error) {
	var n int64
	err := db.Raw(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		schema, table,
	).Scan(&n).Error
	return n > 0, err
}

func existingColumns(db *gorm.DB, schema, table string) (map[string]struct{}, error) {
	var rows []struct{ ColumnName string }
	err := db.Raw(
		"SELECT COLUMN_NAME AS column_name FROM information_schema.columns WHERE table_schema = ? AND table_name = ?",
		schema, table,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{})
	for _, r := range rows {
		m[strings.ToLower(r.ColumnName)] = struct{}{}
	}
	return m, nil
}

// SyncFormTable 根据元数据创建或补齐物理表
func SyncFormTable(form *SysOnlineForm, fields []SysOnlineFormField) error {
	if global.GNA_DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	db := global.GNA_DB
	schema := global.GNA_CONFIG.Database.DbName
	if schema == "" {
		return fmt.Errorf("database.db_name 为空")
	}
	for _, f := range fields {
		if _, skip := reservedCol[strings.ToLower(f.ColumnName)]; skip {
			continue
		}
		if !validateIdent(f.ColumnName) {
			return fmt.Errorf("非法字段名: %s", f.ColumnName)
		}
	}
	exists, err := tableExists(db, schema, form.PhysTableName)
	if err != nil {
		return err
	}
	if !exists {
		sql, err := buildCreateTableSQL(form.PhysTableName, fields)
		if err != nil {
			return err
		}
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("创建表失败: %w", err)
		}
		return nil
	}
	cols, err := existingColumns(db, schema, form.PhysTableName)
	if err != nil {
		return err
	}
	for _, f := range fields {
		if _, skip := reservedCol[strings.ToLower(f.ColumnName)]; skip {
			continue
		}
		if _, ok := cols[strings.ToLower(f.ColumnName)]; ok {
			continue
		}
		colDef, err := columnDDL(f)
		if err != nil {
			return err
		}
		alter := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", form.PhysTableName, colDef)
		if err := db.Exec(alter).Error; err != nil {
			return fmt.Errorf("添加列 %s 失败: %w", f.ColumnName, err)
		}
		cols[strings.ToLower(f.ColumnName)] = struct{}{}
	}
	return nil
}
