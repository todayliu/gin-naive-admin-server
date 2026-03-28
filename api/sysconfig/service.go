package sysconfig

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type _sysConfigService struct{}

var SysConfigService = new(_sysConfigService)

// ConfigEditRequest 编辑系统参数请求体
type ConfigEditRequest struct {
	ID          uint   `json:"id" binding:"required"`
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigValue string `json:"configValue" binding:"required"`
	Remark      string `json:"remark"`
}

// siteDisplayKeys 登录页/壳层标题与页脚：优先使用短键 title、copyright；兼容旧键 site.title、site.copyright
var siteDisplayKeys = []string{"title", "copyright", "site.title", "site.copyright"}

func pickSiteTitle(m map[string]string) string {
	if v := m["title"]; v != "" {
		return v
	}
	return m["site.title"]
}

func pickSiteCopyright(m map[string]string) string {
	if v := m["copyright"]; v != "" {
		return v
	}
	return m["site.copyright"]
}

// SiteDisplay 匿名可访问的站点展示信息（供登录页与壳层标题使用）
// @Summary     站点展示（匿名）
// @Description 返回站点标题、版权等，无需登录。
// @Tags        系统配置
// @Produce     json
// @Success     200 {object} response.Response
// @Router      /config/site-display [get]
func (s *_sysConfigService) SiteDisplay(c *gin.Context) {
	var rows []SysConfig
	if err := global.GNA_DB.Where("config_key IN ?", siteDisplayKeys).Find(&rows).Error; err != nil {
		global.GNA_LOG.Error("读取站点展示配置失败", zap.Error(err))
		response.FailWithMessage("读取失败", c)
		return
	}
	m := make(map[string]string, len(rows))
	for i := range rows {
		m[rows[i].ConfigKey] = rows[i].ConfigValue
	}
	response.OkWithData(gin.H{
		"title":     pickSiteTitle(m),
		"copyright": pickSiteCopyright(m),
	}, c)
}

// List 系统参数列表
// @Summary     系统参数列表
// @Tags        系统配置
// @Produce     json
// @Security    AccessToken
// @Success     200 {object} response.Response
// @Router      /config/list [get]
func (s *_sysConfigService) List(c *gin.Context) {
	var list []SysConfig
	if err := global.GNA_DB.Order("config_key asc").Find(&list).Error; err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(list, c)
}

// ConfigAddRequest 新增系统参数请求体
type ConfigAddRequest struct {
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigValue string `json:"configValue" binding:"required"`
	Remark      string `json:"remark"`
}

// Add 新增一条系统参数（参数键唯一）
// @Summary     新增系统参数
// @Tags        系统配置
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body ConfigAddRequest true "请求体"
// @Success     200 {object} response.Response
// @Router      /config/add [post]
func (s *_sysConfigService) Add(c *gin.Context) {
	var req ConfigAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	var n int64
	global.GNA_DB.Model(&SysConfig{}).Where("config_key = ?", req.ConfigKey).Count(&n)
	if n > 0 {
		response.FailWithMessage("参数键已存在", c)
		return
	}
	row := SysConfig{
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Remark:      req.Remark,
	}
	if err := global.GNA_DB.Create(&row).Error; err != nil {
		global.GNA_LOG.Error("新增参数失败", zap.Error(err))
		response.FailWithMessage("新增失败", c)
		return
	}
	response.Ok(c)
}

// Edit 编辑系统参数
// @Summary     编辑系统参数
// @Tags        系统配置
// @Accept      json
// @Produce     json
// @Security    AccessToken
// @Param       body body ConfigEditRequest true "请求体"
// @Success     200 {object} response.Response
// @Router      /config/edit [put]
func (s *_sysConfigService) Edit(c *gin.Context) {
	var req ConfigEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	err := global.GNA_DB.Model(&SysConfig{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"config_key":   req.ConfigKey,
		"config_value": req.ConfigValue,
		"remark":       req.Remark,
	}).Error
	if err != nil {
		global.GNA_LOG.Error("保存参数失败", zap.Error(err))
		response.FailWithMessage("保存失败", c)
		return
	}
	response.Ok(c)
}

// SeedDefaults 首次库空时写入示例参数
func SeedDefaults(db *gorm.DB) {
	if db == nil {
		return
	}
	var n int64
	db.Model(&SysConfig{}).Count(&n)
	if n > 0 {
		return
	}
	rows := []SysConfig{
		{ConfigKey: "title", ConfigValue: "Gin Naive Admin", Remark: "站点标题"},
		{ConfigKey: "copyright", ConfigValue: "Copyright © 2025", Remark: "页脚版权"},
	}
	_ = db.Create(&rows).Error
}
