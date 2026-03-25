package sysconfig

import "gin-admin-server/global"

// SysConfig 系统参数（键值）
type SysConfig struct {
	global.GNA_MODEL
	ConfigKey   string `gorm:"column:config_key;uniqueIndex;size:128;not null" json:"configKey"`
	ConfigValue string `gorm:"column:config_value;type:text" json:"configValue"`
	Remark      string `gorm:"column:remark;size:255" json:"remark"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}
