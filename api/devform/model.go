package devform

import (
	"gin-admin-server/global"
	"gin-admin-server/model/page_info"
)

// SysOnlineForm 在线表单定义（元数据）
type SysOnlineForm struct {
	global.GNA_MODEL
	PhysTableName string `gorm:"column:table_name;type:varchar(64);uniqueIndex;not null;comment:物理表名" json:"tableName"`
	EntityName    string `gorm:"column:entity_name;type:varchar(64);not null;comment:Go 实体名 PascalCase" json:"entityName"`
	RouteGroup  string `gorm:"column:route_group;type:varchar(64);not null;comment:Gin 路由组段" json:"routeGroup"`
	Description string `gorm:"column:description;type:varchar(255);comment:说明" json:"description"`
	SyncStatus  int    `gorm:"column:sync_status;default:0;comment:0未同步1已同步" json:"syncStatus"`
}

func (SysOnlineForm) TableName() string { return "sys_online_form" }

// SysOnlineFormField 在线表单字段
type SysOnlineFormField struct {
	global.GNA_MODEL
	OnlineFormID uint   `gorm:"column:online_form_id;index;not null;comment:表单ID" json:"onlineFormId"`
	Sort         int    `gorm:"column:sort;default:0;comment:排序" json:"sort"`
	ColumnName   string `gorm:"column:column_name;type:varchar(64);not null;comment:列名snake" json:"columnName"`
	DbType       string `gorm:"column:db_type;type:varchar(32);not null;comment:mysql类型" json:"dbType"`
	Length       int    `gorm:"column:length;default:255;comment:varchar长度" json:"length"`
	DecimalScale int    `gorm:"column:decimal_scale;default:2;comment:decimal小数位" json:"decimalScale"`
	Nullable     bool   `gorm:"column:nullable;default:false" json:"nullable"`
	Comment      string `gorm:"column:comment;type:varchar(255);comment:列注释" json:"comment"`
	ListShow     bool   `gorm:"column:list_show;default:true;comment:列表显示" json:"listShow"`
	FormShow     bool   `gorm:"column:form_show;default:true;comment:表单显示" json:"formShow"`
	IsQuery      bool   `gorm:"column:is_query;default:false;comment:可查询" json:"isQuery"`
	QueryType    string `gorm:"column:query_type;type:varchar(16);default:eq;comment:eq或like" json:"queryType"`
}

func (SysOnlineFormField) TableName() string { return "sys_online_form_field" }

// DevFormPageRequest 表单分页查询
type DevFormPageRequest struct {
	page_info.PageInfo
	TableName   string `form:"tableName"`
	Description string `form:"description"`
}

// DevFormSaveRequest 新增/编辑表单主体
type DevFormSaveRequest struct {
	ID          uint   `json:"id"`
	TableName   string `json:"tableName" binding:"required,max=64"`
	EntityName  string `json:"entityName" binding:"required,max=64"`
	RouteGroup  string `json:"routeGroup" binding:"required,max=64"`
	Description string `json:"description"`
}

// DevFormFieldSaveRequest 保存字段
type DevFormFieldSaveRequest struct {
	ID           uint   `json:"id"`
	OnlineFormID uint   `json:"onlineFormId" binding:"required"`
	Sort         int    `json:"sort"`
	ColumnName   string `json:"columnName" binding:"required,max=64"`
	DbType       string `json:"dbType" binding:"required"`
	Length       int    `json:"length"`
	DecimalScale int    `json:"decimalScale"`
	Nullable     bool   `json:"nullable"`
	Comment      string `json:"comment"`
	ListShow     bool   `json:"listShow"`
	FormShow     bool   `json:"formShow"`
	IsQuery      bool   `json:"isQuery"`
	QueryType    string `json:"queryType"`
}
