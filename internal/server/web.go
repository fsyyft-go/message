// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package server 实现服务器层，负责 HTTP 服务器的配置和启动。
package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/api/sms"
	"github.com/fsyyft-go/sms-bridge/internal/config"
	"github.com/fsyyft-go/sms-bridge/pkg/kratos/middleware/basicauth"
	"github.com/fsyyft-go/sms-bridge/pkg/kratos/middleware/validate"
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
	// 只有 ListSms 操作需要认证，返回 true 表示要应用认证中间件。
	return operation == OperationListSms
}

// NewWebServer 创建 Web 服务器实例。
//
// 参数：
//   - logger log.Logger：日志记录器。
//   - cfg *config.Config：应用配置。
//   - smsService sms.SmsServiceHTTPServer：短信服务 HTTP 接口实例。
//   - auth Authenticator：认证器实例。
//
// 返回值：
//   - WebServer：Web 服务器接口实例。
//   - func()：清理函数，用于资源释放。
//   - error：初始化过程中可能发生的错误。
func NewWebServer(logger log.Logger, cfg *config.Config, smsService sms.SmsServiceHTTPServer, auth Authenticator) (WebServer, func(), error) {
	var err error

	// 创建带有领域驱动设计和模块标记的日志记录器。
	l := logger.WithField("ddd", "server").WithField("module", "web")

	// 初始化 webServer 结构体。
	webServer := &webServer{
		logger: l,
		cfg:    cfg,
	}

	// 创建带基本认证的选择器中间件。
	var authMiddleware middleware.Middleware
	if cfg.Http.Basicauth.Enabled {
		// 如果启用了基本认证，创建带有认证器的中间件。
		authMiddleware = selector.Server(
			basicauth.Server(
				basicauth.WithValidator(auth.Authenticate),    // 使用认证器实例的 Authenticate 方法。
				basicauth.WithRealm(cfg.Http.Basicauth.Realm), // 设置认证域，会显示在浏览器的认证对话框中。
			),
		).Match(needAuthMatcher).Build()
	} else {
		// 不启用认证时，使用空的选择器中间件。
		l.Warn("基本认证未启用，所有请求将不进行认证")
		authMiddleware = selector.Server().Match(func(ctx context.Context, operation string) bool {
			return false // 不匹配任何路径，相当于不应用认证。
		}).Build()
	}

	// 创建 HTTP 服务器，配置中间件链。
	server := http.NewServer(http.Middleware(
		recovery.Recovery(), // 添加恢复中间件，处理 panic。
		validate.Validator(validate.WithValidateCallback(webServer.validateCallback)), // 添加请求验证中间件。
		authMiddleware, // 添加带选择器的认证中间件。
	))
	// 注册短信服务的 HTTP 处理函数。
	sms.RegisterSmsServiceHTTPServer(server, smsService)

	// 初始化 Gin 引擎，并配置默认中间件。
	webServer.engine = gin.Default()
	// 将 Kratos HTTP 服务解析到 Gin 引擎中。
	bridge_kratos_http.Parse(server, webServer.engine)

	// 定义清理函数，用于资源释放。
	var cleanup = func() {}

	// 返回 Web 服务器实例、清理函数和错误。
	return webServer, cleanup, err
}

// validateCallback 处理请求验证失败的回调函数。
// 记录请求和验证错误，并返回标准化的错误响应。
//
// 参数：
//   - ctx context.Context：上下文。
//   - req interface{}：原始请求。
//   - errValidate error：验证过程中产生的错误。
//
// 返回值：
//   - interface{}：处理后的请求（本实现中返回 nil）。
//   - error：格式化后的错误信息。
func (s *webServer) validateCallback(ctx context.Context, req interface{}, errValidate error) (interface{}, error) {
	// 记录请求和验证错误信息。
	s.logger.WithField("req", req).WithField("errValidate", errValidate).Info("validateCallback")
	// 返回标准化的错误响应。
	return nil, errors.BadRequest("VALIDATOR", "请求参数错误，详见日志")
}

// Start 实现启动 Web 服务器的功能。
// 使用 Gin 引擎监听指定端口。
//
// 返回值：
//   - error：启动过程中可能发生的错误。
func (s *webServer) Start() error {
	// 使用 Gin 引擎在配置的端口上启动 HTTP 服务。
	return s.engine.Run(fmt.Sprintf(":%d", s.cfg.Http.Port))
}
