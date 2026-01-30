package config

type SystemConfig struct {
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}
