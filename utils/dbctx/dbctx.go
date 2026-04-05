package dbctx

import (
	"gin-admin-server/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Use 返回带当前请求 context 的 DB 会话，与 JWT 中间件注入的操作人 ID 联动，用于审计字段自动填充
func Use(c *gin.Context) *gorm.DB {
	if c == nil || global.GNA_DB == nil {
		return global.GNA_DB
	}
	return global.GNA_DB.WithContext(c.Request.Context())
}
