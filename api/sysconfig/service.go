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

type configEditRequest struct {
	ID          uint   `json:"id" binding:"required"`
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigValue string `json:"configValue"`
	Remark      string `json:"remark"`
}

// SiteDisplay 匿名可访问的站点展示信息（仅白名单配置键，供登录页与壳层标题使用）
func (s *_sysConfigService) SiteDisplay(c *gin.Context) {
	keys := []string{"site.title", "site.copyright"}
	var rows []SysConfig
	if err := global.GNA_DB.Where("config_key IN ?", keys).Find(&rows).Error; err != nil {
		global.GNA_LOG.Error("读取站点展示配置失败", zap.Error(err))
		response.FailWithMessage("读取失败", c)
		return
	}
	m := make(map[string]string, len(rows))
	for i := range rows {
		m[rows[i].ConfigKey] = rows[i].ConfigValue
	}
	response.OkWithData(gin.H{
		"title":     m["site.title"],
		"copyright": m["site.copyright"],
	}, c)
}

func (s *_sysConfigService) List(c *gin.Context) {
	var list []SysConfig
	if err := global.GNA_DB.Order("config_key asc").Find(&list).Error; err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(list, c)
}

func (s *_sysConfigService) Edit(c *gin.Context) {
	var req configEditRequest
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
		{ConfigKey: "site.title", ConfigValue: "Gin Naive Admin", Remark: "站点标题"},
		{ConfigKey: "site.copyright", ConfigValue: "", Remark: "页脚版权"},
	}
	_ = db.Create(&rows).Error
}
