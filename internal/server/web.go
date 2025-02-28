// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/api/sms"
	"github.com/fsyyft-go/sms-bridge/internal/config"
	"github.com/fsyyft-go/sms-bridge/internal/service"
	bridge_kratos_http "github.com/fsyyft-go/sms-bridge/pkg/kratos/transport/http"
)

type (
	WebServer struct {
		logger log.Logger
		cfg    *config.Config
		engine *gin.Engine
	}
)

func NewWebServer(logger log.Logger, cfg *config.Config, smsService *service.SmsService) (*WebServer, func(), error) {
	var err error
	webServer := &WebServer{
		logger: logger,
		cfg:    cfg,
	}

	server := http.NewServer(http.Middleware(
		recovery.Recovery(),
		validate.Validator(),
	))
	sms.RegisterSmsServiceHTTPServer(server, smsService)

	webServer.engine = gin.Default()
	bridge_kratos_http.Parse(server, webServer.engine)

	var cleanup = func() {}

	return webServer, cleanup, err
}

func (s *WebServer) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.cfg.Http.Port))
}
