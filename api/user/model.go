package user

import (
	"gin-admin-server/global"
	"gin-admin-server/utils/time_util"
)

type User struct {
	global.GNA_MODEL
	UUID          string              `gorm:"column:uuid;type:varchar(100);unique;not null;comment:uuid" json:"uuid"`
	Account       string              `gorm:"column:account;type:varchar(100);unique;not null;comment:用户账号" json:"account"`
	Password      string              `gorm:"column:password;type:varchar(200);not null;comment:用户密码" json:"-"`
	UName         string              `gorm:"column:u_name;type:varchar(100);not null;comment:用户名称" json:"uName"`
	UNickname     string              `gorm:"column:u_nickname;type:varchar(50);comment:用户昵称" json:"uNickname"`
	UMobile       string              `gorm:"column:u_mobile;type:varchar(11);not null;comment:用户手机号码" json:"uMobile"`
	UEmail        string              `gorm:"column:u_email;type:varchar(50);comment:用户邮箱" json:"uEmail"`
	UAvatar       string              `gorm:"column:u_avatar;comment:用户头像" json:"uAvatar"`
	LastLoginTime time_util.LocalTime `gorm:"column:last_login_time;comment:账号最后一次登录时间" json:"lastLoginTime"`
}

func (User) TableName() string {
	return "users"
}
