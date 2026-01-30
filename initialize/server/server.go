package server

import (
	"fmt"
	"gin-admin-server/global"
	"gin-admin-server/initialize/router"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func initServer(port string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func RunServer() {
	Router := router.InitRoute()

	port := fmt.Sprintf(":%d", global.GNA_CONFIG.System.Port)
	s := initServer(port, Router)

	time.Sleep(10 * time.Microsecond)
	global.GNA_LOG.Info("server run success on ", zap.String("address", port))

	fmt.Printf(`
		欢迎使用 gin-naive-admin
		当前版本:v1.0.0
		服务启动成功:http://127.0.0.1%s
		默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	`, port, port)
	global.GNA_LOG.Error(s.ListenAndServe().Error())
}
