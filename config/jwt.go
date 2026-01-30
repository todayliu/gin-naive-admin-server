package config

type JwtConfig struct {
	SecretKey     string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	ExpiresTime   int64  `mapstructure:"expires_time" json:"expires_time" yaml:"expires_time"`
	RefreshTime   int64  `mapstructure:"refresh_time" json:"refresh_time" yaml:"refresh_time"`
	Issuer        string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
	SigningMethod string `mapstructure:"signing_method" json:"signing_method" yaml:"signing_method"`
}
