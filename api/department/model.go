// Package department 部门模块 API
package department

import (
	"gin-admin-server/global"
)

// SysDepartment 部门表
type SysDepartment struct {
	global.GNA_MODEL
	ParentId uint   `gorm:"column:parent_id;not null;default:0;comment:父部门ID" json:"parentId"`
	Name     string `gorm:"column:name;not null;comment:部门名称" json:"name" binding:"required" message:"部门名称不能为空"`
	Code     string `gorm:"column:code;comment:部门编码" json:"code"`
	Sort     uint   `gorm:"column:sort;default:0;comment:排序" json:"sort"`
	Status   uint   `gorm:"column:status;default:1;comment:状态（1:启用 0:禁用）" json:"status"`
	Remark   string `gorm:"column:remark;comment:备注" json:"remark"`
	Children []*SysDepartment `gorm:"-" json:"children,omitempty"`
}

func (SysDepartment) TableName() string {
	return "sys_department"
}
