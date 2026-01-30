package zap_log

import (
	"gin-admin-server/global"
	"os"

	"go.uber.org/zap/zapcore"
)

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct{}

// GetWriteSyncer 获取 zapcore.WriteSyncer
func (r *fileRotatelogs) GetWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriter := NewCutter(
		global.GNA_CONFIG.Zap.Director,
		level,
		WithCutterFormat("2006-01-02"),
		WithMaxSize(global.GNA_CONFIG.Zap.MaxSize),
		WithMaxBackups(global.GNA_CONFIG.Zap.MaxBackups),
		WithMaxAge(global.GNA_CONFIG.Zap.MaxAge),
		WithCompress(global.GNA_CONFIG.Zap.Compress),
	)
	if global.GNA_CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}
	return zapcore.AddSync(fileWriter)
}
