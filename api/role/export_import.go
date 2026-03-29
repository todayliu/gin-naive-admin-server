package role

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gin-admin-server/api/dict"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	roleExportMaxRows = 50000
	roleImportMaxRows = 500
	roleImportMaxBytes = 2 << 20 // 2MB
	dictTypeCodeStatus = "status"
)

func roleImportExportCSVHeader() []string {
	return []string{"角色名称", "角色编码", "状态", "备注"}
}

func dictValueToLabelMap(typeCode string) map[string]string {
	var rows []dict.SysDictData
	if err := global.GNA_DB.Where("type_code = ? AND status = ?", typeCode, 1).Order("sort").Find(&rows).Error; err != nil {
		global.GNA_LOG.Warn("加载字典数据失败", zap.String("typeCode", typeCode), zap.Error(err))
		return nil
	}
	m := make(map[string]string, len(rows))
	for _, r := range rows {
		v := strings.TrimSpace(r.Value)
		if v != "" {
			m[v] = strings.TrimSpace(r.Label)
		}
	}
	return m
}

func exportDictLabel(m map[string]string, raw uint) string {
	if m == nil {
		return strconv.FormatUint(uint64(raw), 10)
	}
	key := strconv.FormatUint(uint64(raw), 10)
	if lab, ok := m[key]; ok && lab != "" {
		return lab
	}
	return key
}

func parseUintOrDictLabel(s string, typeCode string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	if v, err := strconv.ParseUint(s, 10, 32); err == nil {
		return v, nil
	}
	var row dict.SysDictData
	err := global.GNA_DB.Where("type_code = ? AND status = ? AND label = ?", typeCode, 1, s).First(&row).Error
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(row.Value), 10, 32)
}

func buildRoleExportQuery(filters *RoleListFilters) *gorm.DB {
	db := global.GNA_DB.Model(&SysRole{})
	if filters.Name != "" {
		db = db.Where("name LIKE ?", "%"+filters.Name+"%")
	}
	if filters.Code != "" {
		db = db.Where("code LIKE ?", "%"+filters.Code+"%")
	}
	if filters.Status != nil {
		db = db.Where("status = ?", *filters.Status)
	}
	return db
}

// ExportRoles 导出当前筛选条件下的角色为 CSV（UTF-8 BOM）
func (r *_roleService) ExportRoles(c *gin.Context) {
	var filters RoleListFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, filters)
		response.FailWithMessage(errMessage, c)
		return
	}
	var list []SysRole
	if err := buildRoleExportQuery(&filters).Order("create_time desc").Limit(roleExportMaxRows).Find(&list).Error; err != nil {
		global.GNA_LOG.Error("导出角色失败", zap.Error(err))
		response.FailWithMessage("导出角色失败", c)
		return
	}
	statusLabelMap := dictValueToLabelMap(dictTypeCodeStatus)

	buf := &bytes.Buffer{}
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(buf)
	if err := w.Write(roleImportExportCSVHeader()); err != nil {
		global.GNA_LOG.Error("写入 CSV 表头失败", zap.Error(err))
		response.FailWithMessage("导出角色失败", c)
		return
	}
	for _, row := range list {
		st := uint(0)
		if row.Status != nil {
			st = *row.Status
		}
		line := []string{
			row.Name,
			row.Code,
			exportDictLabel(statusLabelMap, st),
			row.Description,
		}
		if err := w.Write(line); err != nil {
			global.GNA_LOG.Error("写入 CSV 行失败", zap.Error(err))
			response.FailWithMessage("导出角色失败", c)
			return
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		global.GNA_LOG.Error("刷新 CSV 失败", zap.Error(err))
		response.FailWithMessage("导出角色失败", c)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="roles_export.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// DownloadRoleImportTemplate 仅含表头的 CSV 模板（UTF-8 BOM）
func (r *_roleService) DownloadRoleImportTemplate(c *gin.Context) {
	var buf bytes.Buffer
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(&buf)
	if err := w.Write(roleImportExportCSVHeader()); err != nil {
		global.GNA_LOG.Error("写入导入模板失败", zap.Error(err))
		response.FailWithMessage("生成模板失败", c)
		return
	}
	w.Flush()
	if err := w.Error(); err != nil {
		global.GNA_LOG.Error("刷新导入模板失败", zap.Error(err))
		response.FailWithMessage("生成模板失败", c)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="role_import_template.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

func canonicalRoleImportKey(header string) string {
	h := strings.TrimSpace(header)
	if h == "" {
		return ""
	}
	switch h {
	case "角色名称", "名称":
		return "name"
	case "角色编码", "编码":
		return "code"
	case "状态":
		return "status"
	case "备注", "描述", "说明":
		return "description"
	}
	lower := strings.ToLower(h)
	switch lower {
	case "name", "rolename":
		return "name"
	case "code", "rolecode":
		return "code"
	case "status":
		return "status"
	case "description", "remark", "desc":
		return "description"
	}
	return lower
}

func buildRoleImportColumnIndex(header []string) map[string]int {
	col := make(map[string]int)
	for i, name := range header {
		key := canonicalRoleImportKey(name)
		if key == "" {
			continue
		}
		col[key] = i
	}
	return col
}

// ImportRolesResult 导入结果
type ImportRolesResult struct {
	SuccessCount int      `json:"successCount"`
	SkipCount    int      `json:"skipCount"`
	FailCount    int      `json:"failCount"`
	Errors       []string `json:"errors"`
}

// ImportRoles 从 CSV 批量导入角色（列与导出一致）
func (r *_roleService) ImportRoles(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("请选择要上传的 CSV 文件", c)
		return
	}
	if file.Size > roleImportMaxBytes {
		response.FailWithMessage("文件过大，请不超过 2MB", c)
		return
	}
	f, err := file.Open()
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	defer f.Close()

	body, err := io.ReadAll(io.LimitReader(f, roleImportMaxBytes+1))
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	if len(body) > roleImportMaxBytes {
		response.FailWithMessage("文件过大，请不超过 2MB", c)
		return
	}
	if len(body) >= 3 && body[0] == 0xef && body[1] == 0xbb && body[2] == 0xbf {
		body = body[3:]
	}

	cr := csv.NewReader(bytes.NewReader(body))
	cr.LazyQuotes = true
	records, err := cr.ReadAll()
	if err != nil {
		response.FailWithMessage("CSV 解析失败: "+err.Error(), c)
		return
	}
	if len(records) < 2 {
		response.FailWithMessage("CSV 至少需要表头与一行数据", c)
		return
	}
	if len(records)-1 > roleImportMaxRows {
		response.FailWithMessage(fmt.Sprintf("单次最多导入 %d 行", roleImportMaxRows), c)
		return
	}

	header := records[0]
	col := buildRoleImportColumnIndex(header)
	required := []string{"name", "code", "status"}
	for _, k := range required {
		if _, ok := col[k]; !ok {
			response.FailWithMessage("缺少必填列（需中文或英文表头）: "+k, c)
			return
		}
	}

	result := ImportRolesResult{Errors: make([]string, 0, 8)}
	const maxErr = 30

	for lineIdx, rec := range records[1:] {
		lineNo := lineIdx + 2
		get := func(keys ...string) string {
			for _, k := range keys {
				if i, ok := col[k]; ok && i < len(rec) {
					return strings.TrimSpace(rec[i])
				}
			}
			return ""
		}
		name := get("name")
		code := get("code")
		statusStr := get("status")
		description := get("description")

		if name == "" || code == "" {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 角色名称与编码不能为空", lineNo))
			}
			continue
		}
		if strings.EqualFold(strings.TrimSpace(code), "admin") {
			result.SkipCount++
			continue
		}
		statusVal, errSt := parseUintOrDictLabel(statusStr, dictTypeCodeStatus)
		if errSt != nil {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 状态须为字典中的数字或标签", lineNo))
			}
			continue
		}
		st := uint(statusVal)

		var exist SysRole
		if err := global.GNA_DB.Where("code = ?", code).First(&exist).Error; err == nil {
			result.SkipCount++
			continue
		} else if err != gorm.ErrRecordNotFound {
			global.GNA_LOG.Error("导入检查角色编码失败", zap.Error(err))
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 数据库错误", lineNo))
			}
			continue
		}

		role := SysRole{
			Name:        name,
			Code:        code,
			Description: description,
			Status:      &st,
		}
		if err := global.GNA_DB.Create(&role).Error; err != nil {
			global.GNA_LOG.Error("导入角色失败", zap.Error(err))
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", lineNo, err))
			}
			continue
		}
		result.SuccessCount++
	}

	msg := fmt.Sprintf("成功 %d，跳过(编码已存在或为保留角色) %d，失败 %d", result.SuccessCount, result.SkipCount, result.FailCount)
	response.OkWithDetailed(result, msg, c)
}
