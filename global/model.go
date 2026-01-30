package global

import (
	"time"

	"gorm.io/gorm"
)

type GNA_MODEL struct {
	ID         uint           `gorm:"column:id;primarykey;autoIncrement" json:"id"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;index;comment:软删除时间" json:"-"`
}
