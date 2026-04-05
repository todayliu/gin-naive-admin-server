package user

import "gin-admin-server/global"

// LoadUserRoleCodes 用户拥有的角色 code 列表（如 admin），用于登录会话与个人资料。
func LoadUserRoleCodes(userID uint) ([]string, error) {
	var codes []string
	err := global.GNA_DB.Table("sys_user_role ur").
		Joins("JOIN sys_role r ON r.id = ur.sys_role_id AND r.delete_time IS NULL").
		Where("ur.sys_user_id = ?", userID).
		Pluck("r.code", &codes).Error
	return codes, err
}
