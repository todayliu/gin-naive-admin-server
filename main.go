package main

import (
	"database/sql"
	"gin-admin-server/global"
	"gin-admin-server/initialize/gorm"
	"gin-admin-server/initialize/redis"
	"gin-admin-server/initialize/server"
	"gin-admin-server/initialize/viper"
	"gin-admin-server/initialize/zap_log"

	"go.uber.org/zap"
)

func main() {
	//初始化 viper
	global.GNA_VIPER = viper.Viper()

	//初始化 zap
	global.GNA_LOG = zap_log.InitZap()
	defer func() {
		_ = global.GNA_LOG.Sync() // 生产环境通常忽略 Stdout 的 Sync 报错
	}()
	zap.ReplaceGlobals(global.GNA_LOG)

	//初始化 gorm 连接数据库
	global.GNA_DB = gorm.InitGorm()
	if global.GNA_DB != nil {
		gorm.RegisterTables() // 初始化表
		db, _ := global.GNA_DB.DB()
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				global.GNA_LOG.Error("数据库关闭错误：" + err.Error())
			}
		}(db)
	}

	//初始化 redis
	global.GNA_REDIS = redis.InitRedis()

	//运行服务
	server.RunServer()
}
