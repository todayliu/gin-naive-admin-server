// Package dict 字典模块 API
package dict

import (
	"gin-admin-server/global"
	"gin-admin-server/model/page_info"
)

// SysDictType 字典类型表
type SysDictType struct {
	global.GNA_MODEL
	TypeCode string `gorm:"column:type_code;unique;not null;comment:字典类型编码" json:"typeCode" binding:"required" message:"字典类型编码不能为空"`
	TypeName string `gorm:"column:type_name;not null;comment:字典类型名称" json:"typeName" binding:"required" message:"字典类型名称不能为空"`
	Status   uint   `gorm:"column:status;default:1;comment:状态（1:启用 0:禁用）" json:"status"`
	Remark   string `gorm:"column:remark;comment:备注" json:"remark"`
	Sort     uint   `gorm:"column:sort;default:0;comment:排序" json:"sort"`
}

func (SysDictType) TableName() string {
	return "sys_dict_type"
}

// SysDictData 字典数据表
type SysDictData struct {
	global.GNA_MODEL
	TypeCode string `gorm:"column:type_code;not null;index;comment:字典类型编码" json:"typeCode" binding:"required" message:"字典类型编码不能为空"`
	Label   string `gorm:"column:label;not null;comment:字典标签" json:"label" binding:"required" message:"字典标签不能为空"`
	Value   string `gorm:"column:value;not null;comment:字典值" json:"value" binding:"required" message:"字典值不能为空"`
	Status  uint   `gorm:"column:status;default:1;comment:状态（1:启用 0:禁用）" json:"status"`
	Remark  string `gorm:"column:remark;comment:备注" json:"remark"`
	Sort    uint   `gorm:"column:sort;default:0;comment:排序" json:"sort"`
}

func (SysDictData) TableName() string {
	return "sys_dict_data"
}

// DictTypePageRequest 字典类型分页查询请求
type DictTypePageRequest struct {
	page_info.PageInfo
	TypeCode string `json:"typeCode" form:"typeCode"` // 字典类型编码
	TypeName string `json:"typeName" form:"typeName"` // 字典类型名称
	Status   *uint  `json:"status" form:"status"`     // 状态
}

// DictDataPageRequest 字典数据分页查询请求
type DictDataPageRequest struct {
	page_info.PageInfo
	TypeCode string `json:"typeCode" form:"typeCode"` // 字典类型编码
	Label    string `json:"label" form:"label"`       // 字典标签
	Status   *uint  `json:"status" form:"status"`      // 状态
}
