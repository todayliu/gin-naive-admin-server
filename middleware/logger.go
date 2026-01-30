package middleware

import (
	"bytes"
	"gin-admin-server/global"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// --- 1. 处理请求体 Body ---
		var body []byte
		// 仅针对具有 Body 的方法进行处理，避开文件上传以防内存溢出
		if (c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch) &&
			!strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {

			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err == nil {
				// 【核心】读完后必须重新写回 Request.Body，否则后续 Controller 拿不到数据
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		c.Next()

		// --- 2. 准备日志字段 ---
		latency := time.Since(start)
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		}

		// --- 3. 将 Body 添加到日志中（带长度限制） ---
		if len(body) > 0 {
			// 限制记录长度为 1024 字符，防止大 JSON 撑爆磁盘
			if len(body) > 1024 {
				fields = append(fields, zap.String("body", string(body[:1024])+"... (truncated)"))
			} else {
				fields = append(fields, zap.String("body", string(body)))
			}
		}

		global.GNA_LOG.Info(path, fields...)
	}
}
