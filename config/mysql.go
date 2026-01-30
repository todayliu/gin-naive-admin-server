package config

import "time"

type DatabaseConfig struct {
	Host            string        `mapstructure:"host" json:"host" yaml:"host"`
	Port            string        `mapstructure:"port" json:"port" yaml:"port"`
	User            string        `mapstructure:"user" json:"user" yaml:"user"`
	Password        string        `mapstructure:"password" json:"password" yaml:"password"`
	DbName          string        `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	Config          string        `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	LogMode         string        `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	LogZap          bool          `mapstructure:"log_zap" json:"log_zap" yaml:"log_zap"`
	Engine          string        `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"`
	Prefix          string        `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Singular        bool          `mapstructure:"singular" json:"singular" yaml:"singular"`
}
