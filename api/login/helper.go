package login

import "gin-admin-server/api/user"

func loadUserRoleCodes(userID uint) []string {
	codes, _ := user.LoadUserRoleCodes(userID)
	return codes
}
