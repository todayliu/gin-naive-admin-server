package config

type Server struct {
	System   SystemConfig   `mapstructure:"system"`
	Database DatabaseConfig `mapstructure:"database"`
	Zap      ZapConfig      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis    RedisConfig    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Router   RouterConfig   `mapstructure:"router" json:"router" yaml:"router"`
	Cors     CorsConfig     `mapstructure:"cors" json:"cors" yaml:"cors"`
	Jwt      JwtConfig      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}
