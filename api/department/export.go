package department

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"sort"
	"strings"

	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ExportDepartmentPaths 导出部门层级路径为单列 CSV（表头「部门路径」；每行一条从根到该节点的完整路径，用 / 分隔；UTF-8 BOM，树序：先序遍历）
func (ds *_departmentService) ExportDepartmentPaths(c *gin.Context) {
	var list []SysDepartment
	if err := dbctx.Use(c).Order("parent_id ASC, sort ASC, id ASC").Find(&list).Error; err != nil {
		global.GNA_LOG.Error("导出部门失败", zap.Error(err))
		response.FailWithMessage("导出部门失败", c)
		return
	}
	byID := make(map[uint]*SysDepartment, len(list))
	for i := range list {
		byID[list[i].ID] = &list[i]
	}

	var lines []string
	var dfs func(parentID uint)
	dfs = func(parentID uint) {
		var children []*SysDepartment
		for i := range list {
			if list[i].ParentId == parentID {
				children = append(children, &list[i])
			}
		}
		sort.Slice(children, func(i, j int) bool {
			if children[i].Sort != children[j].Sort {
				return children[i].Sort < children[j].Sort
			}
			return children[i].ID < children[j].ID
		})
		for _, ch := range children {
			lines = append(lines, deptPathString(ch, byID))
			dfs(ch.ID)
		}
	}
	dfs(0)

	buf := &bytes.Buffer{}
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(buf)
	_ = w.Write([]string{"部门路径"})
	for _, line := range lines {
		_ = w.Write([]string{line})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		global.GNA_LOG.Error("写入 CSV 失败", zap.Error(err))
		response.FailWithMessage("导出部门失败", c)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="department_paths.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

func deptPathString(dept *SysDepartment, byID map[uint]*SysDepartment) string {
	var parts []string
	cur := dept
	const maxHop = 1000
	for hop := 0; hop < maxHop && cur != nil; hop++ {
		parts = append([]string{strings.TrimSpace(cur.Name)}, parts...)
		if cur.ParentId == 0 {
			break
		}
		cur = byID[cur.ParentId]
	}
	return strings.Join(parts, "/")
}
