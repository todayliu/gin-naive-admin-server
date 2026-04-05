package gorm

import (
	"gin-admin-server/global"

	"gorm.io/gorm"
)

// registerAuditCallbacks 在 Create/Update 时写入 CreateBy、UpdateBy；软删除前写入 DeleteBy（依赖请求 context 中的操作人 ID）
func registerAuditCallbacks(db *gorm.DB) {
	if db == nil {
		return
	}
	db.Callback().Create().Before("gorm:before_create").Register("gna:audit_create", auditBeforeCreate)
	db.Callback().Update().Before("gorm:before_update").Register("gna:audit_update", auditBeforeUpdate)
	db.Callback().Delete().After("gorm:before_delete").Register("gna:audit_delete_by", auditInjectDeleteBy)
}

func auditBeforeCreate(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}
	uid := global.OperatorUserID(db.Statement.Context)
	if uid == 0 {
		return
	}
	if f := db.Statement.Schema.LookUpField("CreateBy"); f != nil {
		db.Statement.SetColumn(f.DBName, uid, true)
	}
	if f := db.Statement.Schema.LookUpField("UpdateBy"); f != nil {
		db.Statement.SetColumn(f.DBName, uid, true)
	}
}

func auditBeforeUpdate(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}
	uid := global.OperatorUserID(db.Statement.Context)
	if uid == 0 {
		return
	}
	if f := db.Statement.Schema.LookUpField("UpdateBy"); f != nil {
		db.Statement.SetColumn(f.DBName, uid, true)
	}
}

func auditInjectDeleteBy(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil || db.Statement.Unscoped {
		return
	}
	if db.Statement.Schema.LookUpField("DeleteTime") == nil {
		return
	}
	uid := global.OperatorUserID(db.Statement.Context)
	if uid == 0 {
		return
	}
	df := db.Statement.Schema.LookUpField("DeleteBy")
	if df == nil {
		return
	}
	wc, ok := db.Statement.Clauses["WHERE"]
	if !ok {
		return
	}
	sub := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).WithContext(db.Statement.Context)
	sub.Statement.Schema = db.Statement.Schema
	sub.Statement.Table = db.Statement.Table
	sub.Statement.Model = db.Statement.Model
	sub.Statement.Clauses["WHERE"] = wc
	if err := sub.Updates(map[string]interface{}{df.DBName: uid}).Error; err != nil {
		db.AddError(err)
	}
}
