package log

import (
	"gin-admin-server/global"
	"gin-admin-server/model/page_info"
)

// SysLoginLog 登录日志
type SysLoginLog struct {
	global.GNA_MODEL
	UserId  uint   `gorm:"column:user_id;index;comment:用户ID" json:"userId"`
	Account string `gorm:"column:account;type:varchar(100);index;comment:账号" json:"account"`
	IP      string `gorm:"column:ip;type:varchar(64);comment:IP" json:"ip"`
	Status  int    `gorm:"column:status;index;comment:状态 1成功 2失败" json:"status"`
	Msg     string `gorm:"column:msg;type:varchar(500);comment:说明" json:"msg"`
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}

// SysOperLog 操作日志
type SysOperLog struct {
	global.GNA_MODEL
	Title     string `gorm:"column:title;type:varchar(100);comment:模块标题" json:"title"`
	Method    string `gorm:"column:method;type:varchar(16);index;comment:HTTP方法" json:"method"`
	Path      string `gorm:"column:path;type:varchar(512);comment:路径" json:"path"`
	UserId    uint   `gorm:"column:user_id;index;comment:用户ID" json:"userId"`
	Account   string `gorm:"column:account;type:varchar(100);index;comment:账号" json:"account"`
	IP        string `gorm:"column:ip;type:varchar(64);comment:IP" json:"ip"`
	LatencyMs int64  `gorm:"column:latency_ms;comment:耗时毫秒" json:"latencyMs"`
	Status    int    `gorm:"column:status;comment:HTTP状态码" json:"status"`
}

func (SysOperLog) TableName() string {
	return "sys_oper_log"
}

// LoginLogPageRequest 登录日志分页查询
type LoginLogPageRequest struct {
	page_info.PageInfo
	Account   string `form:"account"`
	Status    string `form:"status"`
	BeginTime string `form:"beginTime"`
	EndTime   string `form:"endTime"`
}

// OperLogPageRequest 操作日志分页查询
type OperLogPageRequest struct {
	page_info.PageInfo
	Account string `form:"account"`
	Method  string `form:"method"`
}
