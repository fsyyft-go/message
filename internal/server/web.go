package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

type (
	WebServer struct {
		logger log.Logger
		cfg    *config.Config
		srv    *gin.Engine
		server *http.Server
	}
)

func NewWebServer(logger log.Logger, cfg *config.Config) (*WebServer, func(), error) {
	var err error

	webServer := &WebServer{
		logger: logger,
		cfg:    cfg,
		srv:    gin.Default(),
	}

	// 设置服务器配置
	webServer.server = &http.Server{
		Addr:    cfg.Http.Addr,
		Handler: webServer.srv,
		// 设置读写超时
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	var cleanup = func() {
		if nil != webServer.server {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := webServer.server.Shutdown(ctx); err != nil {
				webServer.logger.Errorf("关闭 HTTP 服务失败：%v", err)
			}
		}
	}

	return webServer, cleanup, err
}

func (s *WebServer) Start() error {
	return s.server.ListenAndServe()
}
