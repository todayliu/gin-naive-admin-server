package global

import (
	"gin-admin-server/utils/time_util"

	"gorm.io/gorm"
)

type GNA_MODEL struct {
	ID         uint                `gorm:"column:id;primarykey;autoIncrement" json:"id"`
	CreateTime time_util.LocalTime `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	CreateBy   uint                `gorm:"column:create_by;default:0;comment:创建人用户ID" json:"createBy"`
	UpdateTime time_util.LocalTime `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	UpdateBy   uint                `gorm:"column:update_by;default:0;comment:更新人用户ID" json:"updateBy"`
	DeleteTime gorm.DeletedAt      `gorm:"column:delete_time;index;comment:软删除时间" json:"-"`
	DeleteBy   uint                `gorm:"column:delete_by;default:0;comment:删除人用户ID" json:"deleteBy"`
	// 展示字段（无 DB 列）：由 FillAuditDisplayNames 按 create_by/update_by 关联 sys_user.u_name 填充
	CreateByName string `gorm:"-" json:"createByName,omitempty"`
	UpdateByName string `gorm:"-" json:"updateByName,omitempty"`
}
