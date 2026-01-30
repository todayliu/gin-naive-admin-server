package user

import (
	"gin-admin-server/global"
)

type User struct {
	global.GNA_MODEL
	Account   string `gorm:"column:account;type:varchar(100);unique;not null;comment:用户账号" json:"account"`
	Password  string `gorm:"column:password;type:varchar(20);not null;comment:用户密码" json:"-"`
	UName     string `gorm:"column:u_name;type:varchar(100);comment:用户名称" json:"uName"`
	UNickname string `gorm:"column:u_nickname;type:varchar(50);comment:用户昵称" json:"uNickname"`
	UMobile   string `gorm:"column:u_mobile;type:varchar(11);not null;comment:用户手机号码" json:"uMobile"`
	UEmail    string `gorm:"column:u_email;type:varchar(50);comment:用户邮箱" json:"uEmail"`
	UAvatar   string `gorm:"column:u_avatar;comment:用户头像" json:"uAvatar"`
}

func (User) TableName() string {
	return "users"
}
