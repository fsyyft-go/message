// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package server 实现服务器层，负责 HTTP 服务器的配置和启动。
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
	"github.com/fsyyft-go/sms-bridge/pkg/kratos/middleware/basicauth"
	bridge_kratos_http "github.com/fsyyft-go/sms-bridge/pkg/kratos/transport/http"
)

type (
	// WebServer 定义了 Web 服务器的接口。
	WebServer interface {
		// Start 启动 Web 服务器。
		// 返回：
		//   - error：启动过程中可能发生的错误。
		Start() error
	}

	// webServer 实现了 WebServer 接口，提供 Web 服务器功能。
	webServer struct {
		// 日志记录器。
		logger log.Logger
		// 应用配置。
		cfg *config.Config
		// Gin 引擎，用于处理 HTTP 请求。
		engine *gin.Engine
	}
)

// 定义需要认证的操作名称
const (
	// OperationListSms 是查询短信列表操作的名称，需要基本认证。
	OperationListSms = "/api.sms.SmsService/ListSms"
)

// needAuthMatcher 匹配需要认证的操作。
// 用于选择器中间件，决定哪些操作需要应用认证中间件。
//
// 参数：
//   - ctx context.Context：上下文。
//   - operation string：操作名称。
//
// 返回值：
//   - bool：如果操作需要认证，返回 true；否则返回 false。
func needAuthMatcher(ctx context.Context, operation string) bool {
	// 只有 ListSms 操作需要认证，返回 true 表示要应用认证中间件
	return operation == OperationListSms
}

// NewWebServer 创建 Web 服务器实例。
//
// 参数：
//   - logger log.Logger：日志记录器。
//   - cfg *config.Config：应用配置。
//   - smsService sms.SmsServiceHTTPServer：短信服务 HTTP 接口实例。
//
// 返回值：
//   - WebServer：Web 服务器接口实例。
//   - func()：清理函数，用于资源释放。
//   - error：初始化过程中可能发生的错误。
func NewWebServer(logger log.Logger, cfg *config.Config, smsService sms.SmsServiceHTTPServer) (WebServer, func(), error) {
	var err error
	webServer := &webServer{
		logger: logger,
		cfg:    cfg,
	}

	// 创建带基本认证的选择器中间件。
	authMiddleware := selector.Server(
		basicauth.Server(
			basicauth.WithValidator(Authenticate),   // 使用我们定义的认证函数。
			basicauth.WithRealm("SMS Bridge Admin"), // 设置认证域，会显示在浏览器的认证对话框中。
		),
	).Match(needAuthMatcher).Build()

	server := http.NewServer(http.Middleware(
		recovery.Recovery(),
		validate.Validator(),
		authMiddleware, // 添加带选择器的认证中间件。
	))
	sms.RegisterSmsServiceHTTPServer(server, smsService)

	webServer.engine = gin.Default()
	bridge_kratos_http.Parse(server, webServer.engine)

	var cleanup = func() {}

	return webServer, cleanup, err
}

// Start 实现启动 Web 服务器的功能。
// 使用 Gin 引擎监听指定端口。
//
// 返回值：
//   - error：启动过程中可能发生的错误。
func (s *webServer) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.cfg.Http.Port))
}
