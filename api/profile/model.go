package profile

// ProfileUpdateRequest 更新基本信息
type ProfileUpdateRequest struct {
	UNickname string `json:"uNickname"`
	UMobile   string `json:"uMobile" binding:"required"`
	UEmail    string `json:"uEmail"`
	UAvatar   string `json:"uAvatar"`
}

// ProfilePasswordRequest 修改密码（前端传入的 password 为二次 SHA256 后的传输形态，与登录一致）
type ProfilePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
	Salt        string `json:"salt" binding:"required"`
	Timestamp   int64  `json:"timestamp" binding:"required"`
}
