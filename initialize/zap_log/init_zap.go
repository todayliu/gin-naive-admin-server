package zap_log

import (
	"fmt"
	"gin-admin-server/global"
	"gin-admin-server/utils"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap 获取 zap.Logger
func InitZap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.GNA_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.GNA_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GNA_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.GNA_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
