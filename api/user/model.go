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
