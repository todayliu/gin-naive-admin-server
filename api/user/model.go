package user

import (
	"gin-admin-server/global"
	"gin-admin-server/model/page_info"
	"gin-admin-server/utils/time_util"
)

type SysUser struct {
	global.GNA_MODEL
	UUID          string              `gorm:"column:uuid;type:varchar(100);unique;not null;comment:uuid" json:"-"`
	Account       string              `gorm:"column:account;type:varchar(100);unique;not null;comment:用户账号" json:"account"`
	Password      string              `gorm:"column:password;type:varchar(200);not null;comment:用户密码" json:"-"`
	UName         string              `gorm:"column:u_name;type:varchar(100);not null;comment:用户名称" json:"uName"`
	UNickname     string              `gorm:"column:u_nickname;type:varchar(50);comment:用户昵称" json:"uNickname"`
	UMobile       string              `gorm:"column:u_mobile;type:varchar(11);not null;comment:用户手机号码" json:"uMobile"`
	UEmail        string              `gorm:"column:u_email;type:varchar(50);comment:用户邮箱" json:"uEmail"`
	UAvatar       string              `gorm:"column:u_avatar;comment:用户头像" json:"uAvatar"`
	Gender        uint                `gorm:"column:gender;comment:用户性别" json:"gender"`
	Status        uint                `gorm:"column:status;comment:用户状态" json:"status"`
	LastLoginTime time_util.LocalTime `gorm:"column:last_login_time;comment:账号最后一次登录时间" json:"lastLoginTime"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

type UserPageRequest struct {
	page_info.PageInfo
}

// UserAddRequest 新增用户请求
type UserAddRequest struct {
	Account   string `json:"account" binding:"required" message:"用户账号不能为空"`
	Password  string `json:"password" binding:"required,min=6" message:"密码不能为空且至少6位"`
	UName     string `json:"uName" binding:"required" message:"用户名称不能为空"`
	UNickname string `json:"uNickname"`
	UMobile   string `json:"uMobile" binding:"required" message:"手机号不能为空"`
	UEmail    string `json:"uEmail"`
	UAvatar   string `json:"uAvatar"`
	Gender    uint   `json:"gender"`
	Status    uint   `json:"status"`
}

// UserEditRequest 编辑用户请求
type UserEditRequest struct {
	ID        uint   `json:"id" binding:"required" message:"用户ID不能为空"`
	Account   string `json:"account" binding:"required" message:"用户账号不能为空"`
	UName     string `json:"uName" binding:"required" message:"用户名称不能为空"`
	UNickname string `json:"uNickname"`
	UMobile   string `json:"uMobile" binding:"required" message:"手机号不能为空"`
	UEmail    string `json:"uEmail"`
	UAvatar   string `json:"uAvatar"`
	Gender    uint   `json:"gender"`
	Status    uint   `json:"status"`
	Password  string `json:"password"` // 可选，传则更新密码
}
