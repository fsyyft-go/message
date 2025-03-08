// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/api/sms"
	"github.com/fsyyft-go/sms-bridge/internal/config"
	"github.com/fsyyft-go/sms-bridge/internal/middleware/basicauth"
	bridge_kratos_http "github.com/fsyyft-go/sms-bridge/pkg/kratos/transport/http"
)

type (
	WebServer interface {
		Start() error
	}

	webServer struct {
		logger log.Logger
		cfg    *config.Config
		engine *gin.Engine
	}
)

// 定义需要认证的操作名称
const (
	OperationListSms = "/api.sms.SmsService/ListSms"
)

// 匹配需要认证的操作
func needAuthMatcher(ctx context.Context, operation string) bool {
	// 只有 ListSms 操作需要认证，返回 true 表示要应用认证中间件
	return operation == OperationListSms
}

func NewWebServer(logger log.Logger, cfg *config.Config, smsService sms.SmsServiceHTTPServer) (WebServer, func(), error) {
	var err error
	webServer := &webServer{
		logger: logger,
		cfg:    cfg,
	}

	// 创建带基本认证的选择器中间件
	authMiddleware := selector.Server(
		basicauth.Server(
			basicauth.WithValidator(Authenticate),   // 使用我们定义的认证函数
			basicauth.WithRealm("SMS Bridge Admin"), // 设置认证域，会显示在浏览器的认证对话框中
		),
	).Match(needAuthMatcher).Build()

	server := http.NewServer(http.Middleware(
		recovery.Recovery(),
		validate.Validator(),
		authMiddleware, // 添加带选择器的认证中间件
	))
	sms.RegisterSmsServiceHTTPServer(server, smsService)

	webServer.engine = gin.Default()
	bridge_kratos_http.Parse(server, webServer.engine)

	var cleanup = func() {}

	return webServer, cleanup, err
}

func (s *webServer) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.cfg.Http.Port))
}
