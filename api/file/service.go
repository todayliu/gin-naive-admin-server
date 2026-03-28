package file

import (
	"fmt"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type _fileService struct{}

var FileService = new(_fileService)

// Upload 上传文件（multipart）
// @Summary     文件上传
// @Description form 字段名 `file`；成功返回 url、path（相对静态目录）。
// @Tags        文件
// @Accept      multipart/form-data
// @Produce     json
// @Security    AccessToken
// @Param       file formData file true "文件"
// @Success     200 {object} response.Response
// @Router      /file/upload [post]
func (s *_fileService) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("请选择文件", c)
		return
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext == "" {
		ext = ".bin"
	}
	name := utils.GenerateUUID() + ext
	baseDir := global.GNA_CONFIG.Router.Path
	if baseDir == "" {
		baseDir = "uploads/file"
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		global.GNA_LOG.Error("创建上传目录失败", zap.Error(err))
		response.FailWithMessage("上传失败", c)
		return
	}
	dst := filepath.Join(baseDir, name)
	if err := c.SaveUploadedFile(fh, dst); err != nil {
		global.GNA_LOG.Error("保存文件失败", zap.Error(err))
		response.FailWithMessage("上传失败", c)
		return
	}
	urlPath := strings.TrimPrefix(baseDir, "/")
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	url := fmt.Sprintf("%s/%s", strings.TrimRight(urlPath, "/"), name)
	response.OkWithData(gin.H{"url": url, "path": url}, c)
}
