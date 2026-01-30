package config

type ZapConfig struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 日志格式
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Director      string `mapstructure:"director" json:"director" yaml:"director"`                   // 日志目录
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 是否显示行号
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 日志级别编码器
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 堆栈跟踪键名
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 是否输出到控制台
	LocalTime     bool   `mapstructure:"local-time" json:"local-time" yaml:"local-time"`             //是否使用本地时间而不是 UTC
	MaxAge        int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 最大保留天数
	MaxSize       int    `mapstructure:"max-size" json:"max-size" yaml:"max-size"`                   // 每个文件最大
	MaxBackups    int    `mapstructure:"max-backups" json:"max-backups" yaml:"max-backups"`          // 保留多少个备份
	Compress      bool   `mapstructure:"compress" json:"compress" yaml:"compress"`                   //是否开启压缩
}
