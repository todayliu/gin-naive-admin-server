package dict

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type _dictService struct{}

var DictService = new(_dictService)

// ========== 字典类型 ==========

// GetDictTypeList 分页查询字典类型列表
// @Summary     字典类型分页列表
// @Tags        字典
// @Produce     json
// @Security    AccessToken
// @Param       pageNo query int true "页码"
// @Param       pageSize query int true "每页条数"
// @Param       typeCode query string false "类型编码模糊"
// @Param       typeName query string false "类型名称模糊"
// @Success     200 {object} response.Response
// @Router      /dict/type/list [get]
func (d *_dictService) GetDictTypeList(c *gin.Context) {
	var list []SysDictType
	var req DictTypePageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	db := global.GNA_DB.Model(&SysDictType{})
	if req.TypeCode != "" {
		db = db.Where("type_code LIKE ?", "%"+req.TypeCode+"%")
	}
	if req.TypeName != "" {
		db = db.Where("type_name LIKE ?", "%"+req.TypeName+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("获取字典类型总数失败", zap.Error(err))
		response.FailWithMessage("获取字典类型总数失败", c)
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.PageNo - 1)
	if err := db.Limit(limit).Offset(offset).Order("sort ASC, create_time DESC").Find(&list).Error; err != nil {
		global.GNA_LOG.Error("获取字典类型列表失败", zap.Error(err))
		response.FailWithMessage("获取字典类型列表失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}, c)
}

// AddDictType 新增字典类型（若存在已软删除的相同编码则恢复）
// @Summary     新增字典类型
// @Tags        字典
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysDictType true "字典类型"
// @Success     200 {object} response.Response
// @Router      /dict/type/add [post]
func (d *_dictService) AddDictType(c *gin.Context) {
	var req SysDictType
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	// 检查是否存在已软删除的相同类型编码，若存在则恢复
	var deleted SysDictType
	err := global.GNA_DB.Unscoped().Where("type_code = ? AND delete_time IS NOT NULL", req.TypeCode).First(&deleted).Error
	if err == nil {
		if err := global.GNA_DB.Unscoped().Model(&SysDictType{}).Where("id = ?", deleted.ID).Updates(map[string]interface{}{
			"delete_time": nil,
			"type_name":   req.TypeName,
			"status":      req.Status,
			"remark":      req.Remark,
			"sort":        req.Sort,
		}).Error; err != nil {
			global.GNA_LOG.Error("恢复字典类型失败", zap.Error(err))
			response.FailWithMessage("添加失败", c)
			return
		}
		response.Ok(c)
		return
	}
	if err != gorm.ErrRecordNotFound {
		global.GNA_LOG.Error("检查字典类型编码失败", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}
	// 检查是否存在未删除的重复
	var exist SysDictType
	if err := global.GNA_DB.Where("type_code = ?", req.TypeCode).First(&exist).Error; err == nil {
		response.FailWithMessage("字典类型编码已存在", c)
		return
	}
	if err != gorm.ErrRecordNotFound {
		global.GNA_LOG.Error("检查字典类型编码失败", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}
	if err := global.GNA_DB.Create(&req).Error; err != nil {
		global.GNA_LOG.Error("添加字典类型失败", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}
	response.Ok(c)
}

// EditDictType 修改字典类型
// @Summary     编辑字典类型
// @Tags        字典
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysDictType true "字典类型"
// @Success     200 {object} response.Response
// @Router      /dict/type/edit [put]
func (d *_dictService) EditDictType(c *gin.Context) {
	var req SysDictType
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	if err := global.GNA_DB.Model(&SysDictType{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"type_name": req.TypeName,
		"status":    req.Status,
		"remark":    req.Remark,
		"sort":      req.Sort,
	}).Error; err != nil {
		global.GNA_LOG.Error("修改字典类型失败", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.Ok(c)
}

// DeleteDictType 删除字典类型（同时软删除其下所有字典数据）
// @Summary     删除字典类型
// @Tags        字典
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "字典类型 ID"
// @Success     200 {object} response.Response
// @Router      /dict/type/delete/{id} [delete]
func (d *_dictService) DeleteDictType(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var dictType SysDictType
	if err := global.GNA_DB.Where("id = ?", id).First(&dictType).Error; err != nil {
		response.FailWithMessage("字典类型不存在", c)
		return
	}
	if err := global.GNA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("type_code = ?", dictType.TypeCode).Delete(&SysDictData{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&SysDictType{}).Error
	}); err != nil {
		global.GNA_LOG.Error("删除字典类型失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok(c)
}

// ========== 字典数据 ==========

// GetDictDataList 分页查询字典数据列表
// @Summary     字典数据分页列表
// @Tags        字典
// @Produce     json
// @Security    AccessToken
// @Param       pageNo query int true "页码"
// @Param       pageSize query int true "每页条数"
// @Param       typeCode query string false "类型编码"
// @Param       label query string false "标签模糊"
// @Success     200 {object} response.Response
// @Router      /dict/data/list [get]
func (d *_dictService) GetDictDataList(c *gin.Context) {
	var list []SysDictData
	var req DictDataPageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	db := global.GNA_DB.Model(&SysDictData{})
	if req.TypeCode != "" {
		db = db.Where("type_code = ?", req.TypeCode)
	}
	if req.Label != "" {
		db = db.Where("label LIKE ?", "%"+req.Label+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("获取字典数据总数失败", zap.Error(err))
		response.FailWithMessage("获取字典数据总数失败", c)
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.PageNo - 1)
	if err := db.Limit(limit).Offset(offset).Order("sort ASC, create_time DESC").Find(&list).Error; err != nil {
		global.GNA_LOG.Error("获取字典数据列表失败", zap.Error(err))
		response.FailWithMessage("获取字典数据列表失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}, c)
}

// AddDictData 新增字典数据（若存在已软删除的相同数据则恢复）
// @Summary     新增字典数据
// @Tags        字典
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysDictData true "字典数据"
// @Success     200 {object} response.Response
// @Router      /dict/data/add [post]
func (d *_dictService) AddDictData(c *gin.Context) {
	var req SysDictData
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	// 检查是否存在已软删除的相同数据，若存在则恢复
	var deleted SysDictData
	err := global.GNA_DB.Unscoped().Where("type_code = ? AND label = ? AND value = ? AND delete_time IS NOT NULL", req.TypeCode, req.Label, req.Value).First(&deleted).Error
	if err == nil {
		if err := global.GNA_DB.Unscoped().Model(&SysDictData{}).Where("id = ?", deleted.ID).Updates(map[string]interface{}{
			"delete_time": nil,
			"status":      req.Status,
			"remark":      req.Remark,
			"sort":        req.Sort,
		}).Error; err != nil {
			global.GNA_LOG.Error("恢复字典数据失败", zap.Error(err))
			response.FailWithMessage("添加失败", c)
			return
		}
		response.Ok(c)
		return
	}
	if err != gorm.ErrRecordNotFound {
		global.GNA_LOG.Error("检查字典数据失败", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}
	// 检查是否存在未删除的重复数据
	var exist SysDictData
	err = global.GNA_DB.Where("type_code = ? AND label = ? AND value = ?", req.TypeCode, req.Label, req.Value).First(&exist).Error
	if err == nil {
		response.FailWithMessage("该字典类型下已存在相同的标签和值", c)
		return
	}
	if err != gorm.ErrRecordNotFound {
		global.GNA_LOG.Error("检查字典数据失败", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}
	if err := global.GNA_DB.Create(&req).Error; err != nil {
		global.GNA_LOG.Error("添加字典数据失败", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}
	response.Ok(c)
}

// EditDictData 修改字典数据
// @Summary     编辑字典数据
// @Tags        字典
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysDictData true "字典数据"
// @Success     200 {object} response.Response
// @Router      /dict/data/edit [put]
func (d *_dictService) EditDictData(c *gin.Context) {
	var req SysDictData
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	if err := global.GNA_DB.Model(&SysDictData{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"type_code": req.TypeCode,
		"label":     req.Label,
		"value":     req.Value,
		"status":    req.Status,
		"remark":    req.Remark,
		"sort":      req.Sort,
	}).Error; err != nil {
		global.GNA_LOG.Error("修改字典数据失败", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.Ok(c)
}

// DeleteDictData 删除字典数据（软删除）
// @Summary     删除字典数据
// @Tags        字典
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "字典数据 ID"
// @Success     200 {object} response.Response
// @Router      /dict/data/delete/{id} [delete]
func (d *_dictService) DeleteDictData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	if err := global.GNA_DB.Where("id = ?", id).Delete(&SysDictData{}).Error; err != nil {
		global.GNA_LOG.Error("删除字典数据失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok(c)
}

// GetDictByType 根据类型编码获取字典数据（用于下拉选择等）
// @Summary     按类型编码取字典项（登录可访问，用于下拉）
// @Tags        字典
// @Produce     json
// @Security    AccessToken
// @Param       typeCode path string true "类型编码"
// @Success     200 {object} response.Response
// @Router      /dict/data/{typeCode} [get]
func (d *_dictService) GetDictByType(c *gin.Context) {
	typeCode := c.Param("typeCode")
	if typeCode == "" {
		response.FailWithMessage("类型编码不能为空", c)
		return
	}
	var list []SysDictData
	if err := global.GNA_DB.Where("type_code = ? AND status = 1", typeCode).Order("sort ASC").Find(&list).Error; err != nil {
		global.GNA_LOG.Error("获取字典数据失败", zap.Error(err))
		response.FailWithMessage("获取字典数据失败", c)
		return
	}
	response.OkWithData(list, c)
}
