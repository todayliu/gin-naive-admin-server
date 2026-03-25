package login

import (
	"gin-admin-server/api/user"
)

type CaptchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength"`
	Salt          string `json:"salt"`
}
type LoginRequest struct {
	Username  string `json:"username" binding:"required" message:"请输入用户账号"`
	Password  string `json:"password" binding:"required"`
	CaptchaId string `json:"captchaId"`
	Captcha   string `json:"captcha" binding:"required,min=4,max=4" message:"请输入4位验证码"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Salt      string `json:"salt"`
	UUID      string `json:"uuid"`
}

type LoginResponse struct {
	UserInfo   user.SysUser `json:"userInfo"`
	Token      string       `json:"token"`
	Roles      []string     `json:"roles"`
	Codes      []string     `json:"codes"`
	ExpiresIn  int64        `json:"expiresIn"` // 秒，与 Jwt 配置一致
}
