package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

const (
	ERROR    = 201
	SUCCESS  = 200
	NOTFOUND = 404
	TOKEN    = 401
)

func Result(code int, data interface{}, message string, c *gin.Context, success bool) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		message,
		success,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, nil, "操作成功", c, true)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c, true)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c, true)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c, true)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c, false)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c, false)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c, false)
}

func FailNotFound(data interface{}, message string, c *gin.Context) {
	Result(NOTFOUND, data, message, c, false)
}
func FailWithMessageByToken(message string, c *gin.Context) {
	Result(TOKEN, []interface{}{}, message, c, false)
}
