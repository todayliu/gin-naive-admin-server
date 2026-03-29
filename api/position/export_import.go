package position

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	positionExportMaxRows = 50000
	positionImportMaxRows = 500
	positionImportMaxBytes = 2 << 20
)

func positionCSVHeader() []string {
	return []string{"职务级别名称", "职务级别数值"}
}

func buildPositionExportQuery(f *PositionListFilters) *gorm.DB {
	db := global.GNA_DB.Model(&SysJobLevel{})
	if f.LevelName != "" {
		db = db.Where("level_name LIKE ?", "%"+f.LevelName+"%")
	}
	return db
}

// ExportPositions 导出职务级别 CSV（UTF-8 BOM），筛选与列表一致
func (s *_positionService) ExportPositions(c *gin.Context) {
	var filters PositionListFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, filters)
		response.FailWithMessage(errMessage, c)
		return
	}
	var list []SysJobLevel
	if err := buildPositionExportQuery(&filters).Order("level ASC, id ASC").Limit(positionExportMaxRows).Find(&list).Error; err != nil {
		global.GNA_LOG.Error("导出职务失败", zap.Error(err))
		response.FailWithMessage("导出职务失败", c)
		return
	}
	buf := &bytes.Buffer{}
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(buf)
	if err := w.Write(positionCSVHeader()); err != nil {
		response.FailWithMessage("导出职务失败", c)
		return
	}
	for _, row := range list {
		line := []string{row.LevelName, strconv.FormatUint(uint64(row.Level), 10)}
		if err := w.Write(line); err != nil {
			response.FailWithMessage("导出职务失败", c)
			return
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		response.FailWithMessage("导出职务失败", c)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="position_export.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// DownloadPositionImportTemplate 仅表头
func (s *_positionService) DownloadPositionImportTemplate(c *gin.Context) {
	var buf bytes.Buffer
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(&buf)
	_ = w.Write(positionCSVHeader())
	w.Flush()
	c.Header("Content-Disposition", `attachment; filename="position_import_template.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

func canonicalPositionImportKey(h string) string {
	h = strings.TrimSpace(h)
	switch h {
	case "职务级别名称", "级别名称", "名称":
		return "levelname"
	case "职务级别数值", "级别数值", "级别", "数值":
		return "level"
	}
	switch strings.ToLower(h) {
	case "levelname", "level_name", "name":
		return "levelname"
	case "level":
		return "level"
	}
	return strings.ToLower(h)
}

func buildPositionImportCol(header []string) map[string]int {
	col := make(map[string]int)
	for i, name := range header {
		k := canonicalPositionImportKey(name)
		if k != "" {
			col[k] = i
		}
	}
	return col
}

// ImportPositionsResult 导入结果
type ImportPositionsResult struct {
	SuccessCount int      `json:"successCount"`
	SkipCount    int      `json:"skipCount"`
	FailCount    int      `json:"failCount"`
	Errors       []string `json:"errors"`
}

// ImportPositions 从 CSV 导入
func (s *_positionService) ImportPositions(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("请选择要上传的 CSV 文件", c)
		return
	}
	if file.Size > positionImportMaxBytes {
		response.FailWithMessage("文件过大，请不超过 2MB", c)
		return
	}
	f, err := file.Open()
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	defer f.Close()
	body, err := io.ReadAll(io.LimitReader(f, positionImportMaxBytes+1))
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	if len(body) > positionImportMaxBytes {
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
	if len(records)-1 > positionImportMaxRows {
		response.FailWithMessage(fmt.Sprintf("单次最多导入 %d 行", positionImportMaxRows), c)
		return
	}
	col := buildPositionImportCol(records[0])
	for _, k := range []string{"levelname", "level"} {
		if _, ok := col[k]; !ok {
			response.FailWithMessage("缺少必填列: "+k, c)
			return
		}
	}
	result := ImportPositionsResult{Errors: make([]string, 0, 8)}
	const maxErr = 30
	for lineIdx, rec := range records[1:] {
		lineNo := lineIdx + 2
		get := func(key string) string {
			if i, ok := col[key]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		name := get("levelname")
		levelStr := get("level")
		if name == "" || levelStr == "" {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 名称与级别数值不能为空", lineNo))
			}
			continue
		}
		lv, err := strconv.ParseUint(levelStr, 10, 32)
		if err != nil {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 职务级别数值须为非负整数", lineNo))
			}
			continue
		}
		var exist SysJobLevel
		if err := global.GNA_DB.Where("level_name = ?", name).First(&exist).Error; err == nil {
			result.SkipCount++
			continue
		} else if err != gorm.ErrRecordNotFound {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 数据库错误", lineNo))
			}
			continue
		}
		row := SysJobLevel{LevelName: name, Level: uint(lv)}
		if err := global.GNA_DB.Create(&row).Error; err != nil {
			global.GNA_LOG.Error("导入职务失败", zap.Error(err))
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", lineNo, err))
			}
			continue
		}
		result.SuccessCount++
	}
	msg := fmt.Sprintf("成功 %d，跳过(名称已存在) %d，失败 %d", result.SuccessCount, result.SkipCount, result.FailCount)
	response.OkWithDetailed(result, msg, c)
}
