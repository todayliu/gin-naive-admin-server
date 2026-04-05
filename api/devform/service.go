package devform

import (
	"fmt"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"
	"gin-admin-server/utils/validator"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var entityNameRe = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]{0,62}$`)

type _devformService struct{}

var DevformService = new(_devformService)

// GetFormList 在线表单分页列表
func (s *_devformService) GetFormList(c *gin.Context) {
	var req DevFormPageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	db := dbctx.Use(c).Model(&SysOnlineForm{})
	if req.TableName != "" {
		db = db.Where("table_name LIKE ?", "%"+req.TableName+"%")
	}
	if req.Description != "" {
		db = db.Where("description LIKE ?", "%"+req.Description+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("devform list count", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	var list []SysOnlineForm
	limit := req.PageSize
	offset := req.PageSize * (req.PageNo - 1)
	if err := db.Limit(limit).Offset(offset).Order("id DESC").Find(&list).Error; err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	global.FillAuditDisplayNames(dbctx.Use(c), &list)
	response.OkWithData(response.PageResult{List: list, Total: total, PageNo: req.PageNo, PageSize: req.PageSize}, c)
}

func (s *_devformService) QueryForm(c *gin.Context) {
	id := c.Param("id")
	var form SysOnlineForm
	if err := dbctx.Use(c).First(&form, id).Error; err != nil {
		response.FailWithMessage("记录不存在", c)
		return
	}
	if err := ensureGnaModelFormFields(dbctx.Use(c), form.ID); err != nil {
		global.GNA_LOG.Error("devform ensure gna fields", zap.Error(err))
		response.FailWithMessage("系统字段初始化失败: "+err.Error(), c)
		return
	}
	var fields []SysOnlineFormField
	dbctx.Use(c).Where("online_form_id = ?", form.ID).Order("sort ASC, id ASC").Find(&fields)
	global.FillAuditDisplayNames(dbctx.Use(c), &form)
	global.FillAuditDisplayNames(dbctx.Use(c), &fields)
	response.OkWithData(gin.H{"form": form, "fields": fields}, c)
}

func (s *_devformService) AddForm(c *gin.Context) {
	var req DevFormSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	if err := s.validateFormMeta(c, &req, 0); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := purgeSoftDeletedOnlineFormByTableName(dbctx.Use(c), req.TableName); err != nil {
		global.GNA_LOG.Error("devform purge ghost metadata", zap.Error(err))
		response.FailWithMessage("清理同名历史元数据失败: "+err.Error(), c)
		return
	}
	form := SysOnlineForm{
		PhysTableName: req.TableName,
		EntityName:    req.EntityName,
		RouteGroup:  req.RouteGroup,
		Description: req.Description,
		SyncStatus:  0,
	}
	if err := dbctx.Use(c).Create(&form).Error; err != nil {
		global.GNA_LOG.Error("devform add", zap.Error(err))
		response.FailWithMessage("新增失败: "+err.Error(), c)
		return
	}
	if err := seedGnaModelFormFields(dbctx.Use(c), form.ID); err != nil {
		global.GNA_LOG.Error("devform seed gna fields", zap.Error(err))
		response.FailWithMessage("已建表单但写入默认字段失败，请删除后重试", c)
		return
	}
	response.OkWithData(gin.H{"id": form.ID}, c)
}

func (s *_devformService) EditForm(c *gin.Context) {
	var req DevFormSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	if req.ID == 0 {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var old SysOnlineForm
	if err := dbctx.Use(c).First(&old, req.ID).Error; err != nil {
		response.FailWithMessage("记录不存在", c)
		return
	}
	if old.SyncStatus == 1 && req.TableName != old.PhysTableName {
		response.FailWithMessage("已同步后不允许修改物理表名", c)
		return
	}
	if err := s.validateFormMeta(c, &req, req.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	updates := map[string]interface{}{
		"entity_name": req.EntityName,
		"route_group": req.RouteGroup,
		"description": req.Description,
		"table_name":  req.TableName,
	}
	if err := dbctx.Use(c).Model(&SysOnlineForm{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		response.FailWithMessage("保存失败", c)
		return
	}
	response.Ok(c)
}

func (s *_devformService) validateFormMeta(c *gin.Context, req *DevFormSaveRequest, excludeID uint) error {
	if !validateIdent(req.TableName) {
		return fmt.Errorf("表名需匹配小写字母开头，仅含小写、数字、下划线")
	}
	if !validateIdent(req.RouteGroup) {
		return fmt.Errorf("路由组格式非法")
	}
	if !entityNameRe.MatchString(req.EntityName) {
		return fmt.Errorf("实体类名需 PascalCase 英文字母开头")
	}
	var n int64
	q := dbctx.Use(c).Model(&SysOnlineForm{}).Where("table_name = ?", req.TableName)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return fmt.Errorf("表名已存在")
	}
	return nil
}

func (s *_devformService) DeleteForm(c *gin.Context) {
	id := c.Param("id")
	// 元数据表在 table_name 上有唯一索引；软删仍会占位，导致同名表单无法再次新增，故此处硬删字段行与表单行。
	if err := dbctx.Use(c).Unscoped().Where("online_form_id = ?", id).Delete(&SysOnlineFormField{}).Error; err != nil {
		response.FailWithMessage("删除字段失败", c)
		return
	}
	if err := dbctx.Use(c).Unscoped().Delete(&SysOnlineForm{}, id).Error; err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok(c)
}

func (s *_devformService) SaveField(c *gin.Context) {
	var req DevFormFieldSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	if !validateIdent(req.ColumnName) {
		response.FailWithMessage("列名格式非法", c)
		return
	}
	if _, bad := reservedCol[strings.ToLower(req.ColumnName)]; bad {
		response.FailWithMessage("不能使用保留列名（与 GNA 基类一致，已由系统预置）", c)
		return
	}
	if !validateDBType(req.DbType) {
		response.FailWithMessage("不支持的 db_type", c)
		return
	}
	if req.QueryType != "" && req.QueryType != "eq" && req.QueryType != "like" {
		response.FailWithMessage("query_type 仅支持 eq 或 like", c)
		return
	}
	if req.QueryType == "" {
		req.QueryType = "eq"
	}
	var form SysOnlineForm
	if err := dbctx.Use(c).First(&form, req.OnlineFormID).Error; err != nil {
		response.FailWithMessage("表单不存在", c)
		return
	}

	if req.ID == 0 {
		row := SysOnlineFormField{
			OnlineFormID: req.OnlineFormID,
			Sort:         req.Sort,
			ColumnName:   req.ColumnName,
			DbType:       req.DbType,
			Length:       req.Length,
			DecimalScale: req.DecimalScale,
			Nullable:     req.Nullable,
			Comment:      req.Comment,
			ListShow:     req.ListShow,
			FormShow:     req.FormShow,
			IsQuery:      req.IsQuery,
			QueryType:    req.QueryType,
		}
		if err := dbctx.Use(c).Create(&row).Error; err != nil {
			response.FailWithMessage("新增字段失败", c)
			return
		}
	} else {
		var old SysOnlineFormField
		if err := dbctx.Use(c).Where("id = ? AND online_form_id = ?", req.ID, req.OnlineFormID).First(&old).Error; err != nil {
			response.FailWithMessage("字段不存在", c)
			return
		}
		if _, bad := reservedCol[strings.ToLower(old.ColumnName)]; bad {
			response.FailWithMessage("系统预置的 GNA 基类字段不可修改", c)
			return
		}
		updates := map[string]interface{}{
			"sort":          req.Sort,
			"column_name":   req.ColumnName,
			"db_type":       req.DbType,
			"length":        req.Length,
			"decimal_scale": req.DecimalScale,
			"nullable":      req.Nullable,
			"comment":       req.Comment,
			"list_show":     req.ListShow,
			"form_show":     req.FormShow,
			"is_query":      req.IsQuery,
			"query_type":    req.QueryType,
		}
		if err := dbctx.Use(c).Model(&SysOnlineFormField{}).Where("id = ? AND online_form_id = ?", req.ID, req.OnlineFormID).Updates(updates).Error; err != nil {
			response.FailWithMessage("保存字段失败", c)
			return
		}
	}
	response.Ok(c)
}

func (s *_devformService) DeleteField(c *gin.Context) {
	id := c.Param("id")
	var row SysOnlineFormField
	if err := dbctx.Use(c).First(&row, id).Error; err != nil {
		response.FailWithMessage("记录不存在", c)
		return
	}
	if _, bad := reservedCol[strings.ToLower(row.ColumnName)]; bad {
		response.FailWithMessage("系统预置的 GNA 基类字段不可删除", c)
		return
	}
	if err := dbctx.Use(c).Delete(&SysOnlineFormField{}, id).Error; err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok(c)
}

func (s *_devformService) SyncDB(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.FailWithMessage("id 无效", c)
		return
	}
	var form SysOnlineForm
	if err := dbctx.Use(c).First(&form, id).Error; err != nil {
		response.FailWithMessage("表单不存在", c)
		return
	}
	if err := ensureGnaModelFormFields(dbctx.Use(c), form.ID); err != nil {
		global.GNA_LOG.Error("ensure gna form fields", zap.Uint("form_id", form.ID), zap.Error(err))
		response.FailWithMessage("系统字段初始化失败: "+err.Error(), c)
		return
	}
	var fields []SysOnlineFormField
	dbctx.Use(c).Where("online_form_id = ?", form.ID).Order("sort ASC, id ASC").Find(&fields)
	if len(fields) == 0 {
		response.FailWithMessage("请先添加字段", c)
		return
	}
	if err := SyncFormTable(&form, fields); err != nil {
		global.GNA_LOG.Error("sync db", zap.Error(err))
		response.FailWithMessage("同步失败: "+err.Error(), c)
		return
	}
	dbctx.Use(c).Model(&SysOnlineForm{}).Where("id = ?", form.ID).Updates(map[string]interface{}{
		"sync_status": 1,
	})
	response.OkWithMessage("同步成功", c)
}

func (s *_devformService) DownloadCode(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.FailWithMessage("id 无效", c)
		return
	}
	var form SysOnlineForm
	if err := dbctx.Use(c).First(&form, id).Error; err != nil {
		response.FailWithMessage("表单不存在", c)
		return
	}
	if err := ensureGnaModelFormFields(dbctx.Use(c), form.ID); err != nil {
		global.GNA_LOG.Error("ensure gna form fields", zap.Uint("form_id", form.ID), zap.Error(err))
		response.FailWithMessage("系统字段初始化失败: "+err.Error(), c)
		return
	}
	var fields []SysOnlineFormField
	dbctx.Use(c).Where("online_form_id = ?", form.ID).Order("sort ASC, id ASC").Find(&fields)
	if len(fields) == 0 {
		response.FailWithMessage("请先添加字段", c)
		return
	}
	zb, err := ZipGeneratedCode(&form, fields)
	if err != nil {
		global.GNA_LOG.Error("zip codegen", zap.Error(err))
		response.FailWithMessage("生成失败: "+err.Error(), c)
		return
	}
	fn := fmt.Sprintf("codegen_%s_%d.zip", form.PhysTableName, time.Now().Unix())
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=\""+fn+"\"")
	c.Data(http.StatusOK, "application/zip", zb)
}
