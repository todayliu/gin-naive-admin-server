package user

import (
	"gin-admin-server/api/department"
	"gin-admin-server/global"
)

// ScopedDepartmentIDs 当前用户关联部门及其所有子部门 ID（用于列表数据范围）
func ScopedDepartmentIDs(userID uint) []uint {
	var roots []uint
	_ = global.GNA_DB.Model(&SysUserDepartment{}).Where("sys_user_id = ?", userID).Pluck("sys_department_id", &roots)
	var u SysUser
	if err := global.GNA_DB.Select("department_id").Where("id = ?", userID).First(&u).Error; err == nil && u.DepartmentId > 0 {
		roots = append(roots, u.DepartmentId)
	}
	if len(roots) == 0 {
		return nil
	}
	seen := make(map[uint]struct{})
	for _, id := range roots {
		if id > 0 {
			seen[id] = struct{}{}
		}
	}
	if len(seen) == 0 {
		return nil
	}
	var all []department.SysDepartment
	if err := global.GNA_DB.Select("id", "parent_id").Find(&all).Error; err != nil {
		return nil
	}
	children := make(map[uint][]uint)
	for _, d := range all {
		children[d.ParentId] = append(children[d.ParentId], d.ID)
	}
	var stack []uint
	for id := range seen {
		stack = append(stack, id)
	}
	for len(stack) > 0 {
		n := len(stack) - 1
		id := stack[n]
		stack = stack[:n]
		for _, ch := range children[id] {
			if _, ok := seen[ch]; !ok {
				seen[ch] = struct{}{}
				stack = append(stack, ch)
			}
		}
	}
	out := make([]uint, 0, len(seen))
	for id := range seen {
		out = append(out, id)
	}
	return out
}
