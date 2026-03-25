package login

import (
	"gin-admin-server/permission"
)

func loadUserRoleCodes(userID uint) []string {
	codes, _ := permission.LoadUserRoleCodes(userID)
	return codes
}
