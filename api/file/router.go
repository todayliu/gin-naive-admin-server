package file

import "github.com/gin-gonic/gin"

type _fileRouter struct{}

var FileRouter = new(_fileRouter)

func (r *_fileRouter) InitFileRouter(Router *gin.RouterGroup) {
	Router.POST("file/upload", FileService.Upload)
}
