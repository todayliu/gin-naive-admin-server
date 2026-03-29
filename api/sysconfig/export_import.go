package sysconfig

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gin-admin-server/global"
	"gin-admin-server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	configImportMaxRows  = 2000
	configImportMaxBytes = 2 << 20
)

func sysConfigCSVHeader() []string {
	return []string{"参数键", "参数值", "备注"}
}

// ExportConfigs 导出全部系统参数 CSV（UTF-8 BOM）
func (s *_sysConfigService) ExportConfigs(c *gin.Context) {
	var list []SysConfig
	if err := global.GNA_DB.Order("config_key asc").Find(&list).Error; err != nil {
		global.GNA_LOG.Error("导出参数失败", zap.Error(err))
		response.FailWithMessage("导出参数失败", c)
		return
	}
	buf := &bytes.Buffer{}
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(buf)
	if err := w.Write(sysConfigCSVHeader()); err != nil {
		response.FailWithMessage("导出参数失败", c)
		return
	}
	for _, row := range list {
		line := []string{row.ConfigKey, row.ConfigValue, row.Remark}
		if err := w.Write(line); err != nil {
			response.FailWithMessage("导出参数失败", c)
			return
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		response.FailWithMessage("导出参数失败", c)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="sys_config_export.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// DownloadConfigImportTemplate 仅表头
func (s *_sysConfigService) DownloadConfigImportTemplate(c *gin.Context) {
	var buf bytes.Buffer
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(&buf)
	_ = w.Write(sysConfigCSVHeader())
	w.Flush()
	c.Header("Content-Disposition", `attachment; filename="sys_config_import_template.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

func canonicalConfigImportKey(h string) string {
	h = strings.TrimSpace(h)
	switch h {
	case "参数键", "键", "配置键":
		return "configkey"
	case "参数值", "值", "配置值":
		return "configvalue"
	case "备注", "说明":
		return "remark"
	}
	switch strings.ToLower(h) {
	case "configkey", "config_key", "key":
		return "configkey"
	case "configvalue", "config_value", "value":
		return "configvalue"
	case "remark", "description":
		return "remark"
	}
	return strings.ToLower(h)
}

func buildConfigImportCol(header []string) map[string]int {
	col := make(map[string]int)
	for i, name := range header {
		k := canonicalConfigImportKey(name)
		if k != "" {
			col[k] = i
		}
	}
	return col
}

// ImportConfigsResult 导入结果
type ImportConfigsResult struct {
	SuccessCount int      `json:"successCount"`
	SkipCount    int      `json:"skipCount"`
	FailCount    int      `json:"failCount"`
	Errors       []string `json:"errors"`
}

// ImportConfigs 从 CSV 导入（参数键已存在则跳过）
func (s *_sysConfigService) ImportConfigs(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("请选择要上传的 CSV 文件", c)
		return
	}
	if file.Size > configImportMaxBytes {
		response.FailWithMessage("文件过大，请不超过 2MB", c)
		return
	}
	f, err := file.Open()
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	defer f.Close()
	body, err := io.ReadAll(io.LimitReader(f, configImportMaxBytes+1))
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	if len(body) > configImportMaxBytes {
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
	if len(records)-1 > configImportMaxRows {
		response.FailWithMessage(fmt.Sprintf("单次最多导入 %d 行", configImportMaxRows), c)
		return
	}
	col := buildConfigImportCol(records[0])
	for _, k := range []string{"configkey", "configvalue"} {
		if _, ok := col[k]; !ok {
			response.FailWithMessage("缺少必填列: "+k, c)
			return
		}
	}
	result := ImportConfigsResult{Errors: make([]string, 0, 8)}
	const maxErr = 30
	for lineIdx, rec := range records[1:] {
		lineNo := lineIdx + 2
		get := func(key string) string {
			if i, ok := col[key]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		key := get("configkey")
		val := get("configvalue")
		remark := get("remark")
		if key == "" {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 参数键不能为空", lineNo))
			}
			continue
		}
		var n int64
		global.GNA_DB.Model(&SysConfig{}).Where("config_key = ?", key).Count(&n)
		if n > 0 {
			result.SkipCount++
			continue
		}
		row := SysConfig{ConfigKey: key, ConfigValue: val, Remark: remark}
		if err := global.GNA_DB.Create(&row).Error; err != nil {
			global.GNA_LOG.Error("导入参数失败", zap.Error(err))
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", lineNo, err))
			}
			continue
		}
		result.SuccessCount++
	}
	msg := fmt.Sprintf("成功 %d，跳过(参数键已存在) %d，失败 %d", result.SuccessCount, result.SkipCount, result.FailCount)
	response.OkWithDetailed(result, msg, c)
}
