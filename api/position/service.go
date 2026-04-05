package position

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/dbctx"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
)

type _positionService struct{}

var PositionService = new(_positionService)

// GetPositionList 分页查询职务级别列表（按 level 升序，数值越小级别越高）
// @Summary     职务分页列表
// @Tags        职务
// @Produce     json
// @Security    AccessToken
// @Param       pageNo query int true "页码"
// @Param       pageSize query int true "每页条数"
// @Param       levelName query string false "级别名称模糊"
// @Success     200 {object} response.Response
// @Router      /position/list [get]
func (s *_positionService) GetPositionList(c *gin.Context) {
	var list []SysJobLevel
	var req PositionPageRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, req)
		response.FailWithMessage(errMessage, c)
		return
	}
	db := dbctx.Use(c).Model(&SysJobLevel{})
	if req.LevelName != "" {
		db = db.Where("level_name LIKE ?", "%"+req.LevelName+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("获取职务级别总数失败" + err.Error())
		response.FailWithMessage("获取职务级别总数失败", c)
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.PageNo - 1)
	err = db.Limit(limit).Offset(offset).Order("level ASC, id ASC").Find(&list).Error
	if err != nil {
		global.GNA_LOG.Error("获取职务级别列表失败：" + err.Error())
		response.FailWithMessage("获取职务级别列表失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}, c)
}

// QueryPosition 查询职务级别详情
// @Summary     职务详情
// @Tags        职务
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "职务 ID"
// @Success     200 {object} response.Response
// @Router      /position/query/{id} [get]
func (s *_positionService) QueryPosition(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	var row SysJobLevel
	err := dbctx.Use(c).Where("id = ?", id).First(&row).Error
	if err != nil {
		global.GNA_LOG.Error("查询职务级别失败：" + err.Error())
		response.FailWithMessage("查询职务级别失败", c)
		return
	}
	response.OkWithData(row, c)
}

// AddPosition 新增职务级别
// @Summary     新增职务
// @Tags        职务
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysJobLevel true "职务"
// @Success     200 {object} response.Response
// @Router      /position/add [post]
func (s *_positionService) AddPosition(c *gin.Context) {
	var row SysJobLevel
	err := c.ShouldBindJSON(&row)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, row)
		response.FailWithMessage(errMessage, c)
		return
	}
	err = dbctx.Use(c).Create(&row).Error
	if err != nil {
		global.GNA_LOG.Error("新增职务级别失败：" + err.Error())
		response.FailWithMessage("新增职务级别失败", c)
		return
	}
	response.Ok(c)
}

// EditPosition 修改职务级别
// @Summary     编辑职务
// @Tags        职务
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body SysJobLevel true "职务（含 id）"
// @Success     200 {object} response.Response
// @Router      /position/edit [put]
func (s *_positionService) EditPosition(c *gin.Context) {
	var row SysJobLevel
	err := c.ShouldBindJSON(&row)
	if err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, row)
		response.FailWithMessage(errMessage, c)
		return
	}
	if row.ID == 0 {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	err = dbctx.Use(c).Model(&SysJobLevel{}).Where("id = ?", row.ID).Updates(map[string]interface{}{
		"level_name": row.LevelName,
		"level":      row.Level,
	}).Error
	if err != nil {
		global.GNA_LOG.Error("修改职务级别失败：" + err.Error())
		response.FailWithMessage("修改职务级别失败", c)
		return
	}
	response.Ok(c)
}

// DeletePosition 删除职务级别
// @Summary     删除职务
// @Tags        职务
// @Produce     json
// @Security    AccessToken
// @Param       id path int true "职务 ID"
// @Success     200 {object} response.Response
// @Router      /position/delete/{id} [delete]
func (s *_positionService) DeletePosition(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.FailWithMessage("id 不能为空", c)
		return
	}
	err := dbctx.Use(c).Where("id = ?", id).Delete(&SysJobLevel{}).Error
	if err != nil {
		global.GNA_LOG.Error("删除职务级别失败：" + err.Error())
		response.FailWithMessage("删除职务级别失败", c)
		return
	}
	response.Ok(c)
}
