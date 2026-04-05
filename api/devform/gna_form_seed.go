package devform

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// gnaMetaDefaults 与 global.GNA_MODEL + ddl 固定段一致（在线列表展示 / 元数据）
var gnaMetaDefaults = []SysOnlineFormField{
	{Sort: 1, ColumnName: "id", DbType: "bigint", Length: 0, DecimalScale: 0, Nullable: false, Comment: "主键", ListShow: false, FormShow: false, IsQuery: false, QueryType: "eq"},
	{Sort: 2, ColumnName: "create_by", DbType: "bigint", Length: 0, DecimalScale: 0, Nullable: false, Comment: "创建人用户ID", ListShow: false, FormShow: false, IsQuery: false, QueryType: "eq"},
	{Sort: 3, ColumnName: "update_by", DbType: "bigint", Length: 0, DecimalScale: 0, Nullable: false, Comment: "更新人用户ID", ListShow: false, FormShow: false, IsQuery: false, QueryType: "eq"},
	{Sort: 4, ColumnName: "delete_by", DbType: "bigint", Length: 0, DecimalScale: 0, Nullable: false, Comment: "删除人用户ID", ListShow: false, FormShow: false, IsQuery: false, QueryType: "eq"},
	{Sort: 5, ColumnName: "create_time", DbType: "datetime", Length: 0, DecimalScale: 0, Nullable: true, Comment: "创建时间", ListShow: false, FormShow: false, IsQuery: false, QueryType: "eq"},
	{Sort: 6, ColumnName: "update_time", DbType: "datetime", Length: 0, DecimalScale: 0, Nullable: true, Comment: "更新时间", ListShow: false, FormShow: false, IsQuery: false, QueryType: "eq"},
	{Sort: 7, ColumnName: "delete_time", DbType: "datetime", Length: 0, DecimalScale: 0, Nullable: true, Comment: "软删除时间", ListShow: false, FormShow: false, IsQuery: false, QueryType: "eq"},
}

// ensureGnaModelFormFields 幂等补全：缺哪条插哪条（新建前或未跑过种子的历史表单）
func ensureGnaModelFormFields(db *gorm.DB, onlineFormID uint) error {
	if db == nil || onlineFormID == 0 {
		return nil
	}
	for _, def := range gnaMetaDefaults {
		var n int64
		col := strings.ToLower(def.ColumnName)
		if err := db.Model(&SysOnlineFormField{}).
			Where("online_form_id = ? AND LOWER(column_name) = ?", onlineFormID, col).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			if err := db.Model(&SysOnlineFormField{}).
				Where("online_form_id = ? AND LOWER(column_name) = ?", onlineFormID, col).
				Update("sort", def.Sort).Error; err != nil {
				return err
			}
			continue
		}
		row := def
		row.OnlineFormID = onlineFormID
		if err := db.Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

// seedGnaModelFormFields 新建表单时写入默认字段（等同 ensure）
func seedGnaModelFormFields(db *gorm.DB, onlineFormID uint) error {
	return ensureGnaModelFormFields(db, onlineFormID)
}

// purgeSoftDeletedOnlineFormByTableName 释放 table_name 唯一键：历史软删残留在库中会导致无法同名再建。
func purgeSoftDeletedOnlineFormByTableName(db *gorm.DB, tableName string) error {
	if db == nil || tableName == "" {
		return nil
	}
	var ghost SysOnlineForm
	err := db.Unscoped().Where("table_name = ?", tableName).First(&ghost).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if !ghost.DeleteTime.Valid {
		return nil
	}
	if err := db.Unscoped().Where("online_form_id = ?", ghost.ID).Delete(&SysOnlineFormField{}).Error; err != nil {
		return err
	}
	return db.Unscoped().Delete(&SysOnlineForm{}, ghost.ID).Error
}
