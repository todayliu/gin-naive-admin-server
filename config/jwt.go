package config

type JwtConfig struct {
	SecretKey   string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	ExpiresTime int64  `mapstructure:"expires-time" json:"expires-time" yaml:"expires-time"`
	RefreshTime int64  `mapstructure:"refresh-time" json:"refresh-time" yaml:"refresh-time"`
	BufferTime  int64  `mapstructure:"buffer-time" json:"buffer-time" yaml:"buffer-time"`
	Issuer      string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}
