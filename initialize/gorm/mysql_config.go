package gorm

import (
	"gin-admin-server/config"

	"gorm.io/driver/mysql"
)

var Mysql = new(_mysql)

type _mysql struct{}

func (m *_mysql) InitDbConfig(dbConfig *config.DatabaseConfig) mysql.Config {
	dbConf := mysql.Config{
		DSN:                       dbConfig.User + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.DbName + "?" + dbConfig.Config, // DSN data source name
		DefaultStringSize:         191,                                                                                                                                      // string 类型字段的默认长度
		SkipInitializeWithVersion: false,                                                                                                                                    // 根据版本自动配置
	}

	return dbConf
}
