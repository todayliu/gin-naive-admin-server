package config

type RouterConfig struct {
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
}
