package zap_log

import (
	"gin-admin-server/global"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(global.GNA_CONFIG.Zap.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
}

// TransportLevel 根据字符串转化为 zapcore.Level
func TransportLevel() zapcore.Level {
	global.GNA_CONFIG.Zap.Level = strings.ToLower(global.GNA_CONFIG.Zap.Level)
	switch global.GNA_CONFIG.Zap.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
