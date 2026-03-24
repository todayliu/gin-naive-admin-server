package gorm

import (
	"gin-admin-server/api/department"
	"gin-admin-server/api/dict"
	"gin-admin-server/api/position"
	"gin-admin-server/api/menu"
	"gin-admin-server/api/role"
	"gin-admin-server/api/user"
	"gin-admin-server/config"
	"gin-admin-server/global"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDsn(dbConfig *config.DatabaseConfig) string {
	return dbConfig.User + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.DbName + "?" + dbConfig.Config
}

func InitGorm() *gorm.DB {
	var dbConfig = global.GNA_CONFIG.Database
	if dbConfig.DbName == "" {
		return nil
	}

	db, err := gorm.Open(mysql.New(Mysql.InitDbConfig(&dbConfig)), Gorm.GormConfig(dbConfig.Prefix, dbConfig.Singular))
	if err != nil {
		global.GNA_LOG.Error("数据库连接失败，err:" + err.Error())
		return nil
	}
	global.GNA_LOG.Info("数据库连接成功")
	db.InstanceSet("gorm:table_options", "ENGINE="+dbConfig.Engine)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	return db
}

// RegisterTables 注册数据库表
func RegisterTables() {
	db := global.GNA_DB
	err := db.AutoMigrate(
		// 系统模块表
		user.SysUser{},
		user.SysUserDepartment{},
		user.SysUserJobLevel{},
		menu.SysMenu{},
		role.SysRole{},
		department.SysDepartment{},
		dict.SysDictType{},
		dict.SysDictData{},
		position.SysJobLevel{},
	)
	if err != nil {
		global.GNA_LOG.Error("注册数据表失败", zap.Error(err))
		os.Exit(0)
	}
	global.GNA_LOG.Info("注册数据表成功")
}
