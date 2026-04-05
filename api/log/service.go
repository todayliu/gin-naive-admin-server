package log

import (
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils/validator"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type _logService struct{}

var LogService = new(_logService)

// SaveLoginLogAsync 异步写入登录日志
func SaveLoginLogAsync(userId uint, account, ip string, status int, msg string) {
	if global.GNA_DB == nil {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				global.GNA_LOG.Error("SaveLoginLogAsync panic", zap.Any("recover", r))
			}
		}()
		row := SysLoginLog{
			UserId:  userId,
			Account: account,
			IP:      ip,
			Status:  status,
			Msg:     msg,
		}
		if err := global.GNA_DB.Create(&row).Error; err != nil {
			global.GNA_LOG.Error("写入登录日志失败", zap.Error(err))
		}
	}()
}

// SaveOperLogAsync 异步写入操作日志
func SaveOperLogAsync(title, method, path string, userId uint, account, ip string, latencyMs int64, status int) {
	if global.GNA_DB == nil {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				global.GNA_LOG.Error("SaveOperLogAsync panic", zap.Any("recover", r))
			}
		}()
		row := SysOperLog{
			Title:     title,
			Method:    method,
			Path:      path,
			UserId:    userId,
			Account:   account,
			IP:        ip,
			LatencyMs: latencyMs,
			Status:    status,
		}
		if err := global.GNA_DB.Create(&row).Error; err != nil {
			global.GNA_LOG.Error("写入操作日志失败", zap.Error(err))
		}
	}()
}

// GetLoginLogList 登录日志分页列表
// @Summary     登录日志列表
// @Tags        日志
// @Produce     json
// @Security    AccessToken
// @Param       pageNo query int true "页码"
// @Param       pageSize query int true "每页条数"
// @Param       account query string false "账号模糊"
// @Param       status query int false "状态"
// @Param       beginTime query string false "开始时间"
// @Param       endTime query string false "结束时间"
// @Success     200 {object} response.Response
// @Router      /log/login/list [get]
func (s *_logService) GetLoginLogList(c *gin.Context) {
	var req LoginLogPageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	db := global.GNA_DB.Model(&SysLoginLog{})
	if req.Account != "" {
		db = db.Where("account LIKE ?", "%"+req.Account+"%")
	}
	if req.Status != "" {
		if st, err := strconv.Atoi(req.Status); err == nil {
			db = db.Where("status = ?", st)
		}
	}
	db = appendLoginLogTimeFilter(db, req.BeginTime, req.EndTime)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("统计登录日志失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.PageNo - 1)
	var list []SysLoginLog
	if err := db.Limit(limit).Offset(offset).Order("create_time DESC").Find(&list).Error; err != nil {
		global.GNA_LOG.Error("查询登录日志失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}, c)
}

// GetOperLogList 操作日志分页列表
// @Summary     操作日志列表
// @Tags        日志
// @Produce     json
// @Security    AccessToken
// @Param       pageNo query int true "页码"
// @Param       pageSize query int true "每页条数"
// @Param       account query string false "账号模糊"
// @Param       method query string false "HTTP 方法"
// @Success     200 {object} response.Response
// @Router      /log/oper/list [get]
func (s *_logService) GetOperLogList(c *gin.Context) {
	var req OperLogPageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(validator.GetValidatorErrorMessage(err, req), c)
		return
	}
	db := global.GNA_DB.Model(&SysOperLog{})
	if req.Account != "" {
		db = db.Where("account LIKE ?", "%"+req.Account+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.GNA_LOG.Error("统计操作日志失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	limit := req.PageSize
	offset := req.PageSize * (req.PageNo - 1)
	var list []SysOperLog
	if err := db.Limit(limit).Offset(offset).Order("create_time DESC").Find(&list).Error; err != nil {
		global.GNA_LOG.Error("查询操作日志失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}, c)
}

func appendLoginLogTimeFilter(db *gorm.DB, begin, end string) *gorm.DB {
	if begin != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", begin, time.Local); err == nil {
			db = db.Where("create_time >= ?", t)
		} else if t, err := time.Parse(time.RFC3339, begin); err == nil {
			db = db.Where("create_time >= ?", t)
		} else if t, err := time.ParseInLocation("2006-01-02", begin, time.Local); err == nil {
			db = db.Where("create_time >= ?", t)
		}
	}
	if end != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local); err == nil {
			db = db.Where("create_time <= ?", t)
		} else if t, err := time.Parse(time.RFC3339, end); err == nil {
			db = db.Where("create_time <= ?", t)
		} else if t, err := time.ParseInLocation("2006-01-02", end, time.Local); err == nil {
			eod := t.Add(24*time.Hour - time.Nanosecond)
			db = db.Where("create_time <= ?", eod)
		}
	}
	return db
}
