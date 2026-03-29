// Package position 职务管理（职务级别）模块 API
package position

import (
	"gin-admin-server/global"
	"gin-admin-server/model/page_info"
)

// SysJobLevel 职务级别表（表名保持 sys_job_level，与历史数据兼容）
type SysJobLevel struct {
	global.GNA_MODEL
	LevelName string `gorm:"column:level_name;type:varchar(100);not null;comment:职务级别名称" json:"levelName" binding:"required" message:"职务级别名称不能为空"`
	Level     uint   `gorm:"column:level;not null;comment:职务级别数值（越小级别越高）" json:"level" binding:"required" message:"职务级别不能为空"`
}

func (SysJobLevel) TableName() string {
	return "sys_job_level"
}

// PositionListFilters 职务列表/导出共用的查询条件（不含分页）
type PositionListFilters struct {
	LevelName string `form:"levelName"` // 按名称模糊查询
}

// PositionPageRequest 职务级别分页查询请求
type PositionPageRequest struct {
	page_info.PageInfo
	PositionListFilters
}
