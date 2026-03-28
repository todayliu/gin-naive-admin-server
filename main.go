// Package main Gin Naive Admin 服务端入口。
//
// @title Gin Naive Admin API
// @version 1.0
// @description 后台管理 REST 接口。统一响应结构体 `response.Response`：`code` 200 成功、401 Token 无效或过期、201 业务失败；`data` 为具体载荷。除「认证」「系统/健康」「站点展示」外，均需请求头 `AccessToken`。
// @termsOfService https://swagger.io/terms/
//
// @contact.name API Support
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host 127.0.0.1:8080
// @BasePath /api
//
// @securityDefinitions.apikey AccessToken
// @in header
// @name AccessToken
// @description 登录成功后返回的 JWT，请求头字段名为 AccessToken（与前端一致）
package main

import (
	"database/sql"
	"gin-admin-server/global"
	"gin-admin-server/initialize/gorm"
	"gin-admin-server/initialize/redis"
	"gin-admin-server/initialize/server"
	"gin-admin-server/initialize/viper"
	"gin-admin-server/initialize/zap_log"
	"gin-admin-server/api/sysconfig"
	"gin-admin-server/permission"
	"gin-admin-server/utils/validator"

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
	//初始化验证规则规则
	validator.InitValidator()
	//初始化 gorm 连接数据库
	global.GNA_DB = gorm.InitGorm()
	if global.GNA_DB != nil {
		gorm.RegisterTables() // 初始化表
		permission.SeedMenuButtonPermsIfNeeded(global.GNA_DB)
		permission.ReparentAPIPermissionButtons(global.GNA_DB)
		sysconfig.SeedDefaults(global.GNA_DB)
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
