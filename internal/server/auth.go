// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package server 实现服务器层，负责 HTTP 服务器的配置和启动。
package server

import (
	"context"
	"fmt"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/message/internal/config"
)

type (
	// Authenticator 定义了身份验证器的接口。
	Authenticator interface {
		// Authenticate 验证用户名和密码是否正确。
		// 用于基本认证中间件的验证函数。
		//
		// 参数：
		//   - ctx context.Context：上下文。
		//   - username string：用户名。
		//   - password string：密码。
		//
		// 返回值：
		//   - bool：如果认证成功，返回 true；否则返回 false。
		Authenticate(ctx context.Context, username, password string) bool
	}

	// authenticator 实现了 Authenticator 接口，提供身份验证功能。
	authenticator struct {
		// 日志记录器。
		logger log.Logger
		// 应用配置。
		cfg *config.Config
	}
)

// NewAuthenticator 创建一个新的身份验证器实例。
//
// 参数：
//   - logger log.Logger：日志记录器。
//   - cfg *config.Config：应用配置。
//
// 返回值：
//   - Authenticator：身份验证器接口实例。
//   - func()：清理函数，用于资源释放。
//   - error：初始化过程中可能发生的错误。
func NewAuthenticator(logger log.Logger, cfg *config.Config) (Authenticator, func(), error) {
	// 验证必要参数。
	if nil == logger {
		return nil, nil, fmt.Errorf("创建认证器失败：日志记录器不能为空")
	}

	if nil == cfg {
		return nil, nil, fmt.Errorf("创建认证器失败：配置不能为空")
	}

	l := logger.WithField("ddd", "server").WithField("module", "auth")

	// 检查基本认证是否开启。
	if nil != cfg.Http && nil != cfg.Http.Basicauth && !cfg.Http.Basicauth.Enabled {
		l.Warn("基本认证未开启，API 接口将不进行认证检查")
	}

	// 创建认证器实例。
	auth := &authenticator{
		logger: l,
		cfg:    cfg,
	}

	// 创建空的清理函数。
	cleanup := func() {
		// 当前没有需要清理的资源。
	}

	return auth, cleanup, nil
}

// Authenticate 验证用户名和密码是否正确。
// 用于基本认证中间件的验证函数。
//
// 参数：
//   - ctx context.Context：上下文。
//   - username string：用户名。
//   - password string：密码。
//
// 返回值：
//   - bool：如果认证成功，返回 true；否则返回 false。
func (a *authenticator) Authenticate(ctx context.Context, username, password string) bool {
	// 配置不完整的情况。
	if nil == a.cfg {
		a.logger.Error("认证失败：配置对象为空")
		return false
	}

	if nil == a.cfg.Http {
		a.logger.Error("认证失败：HTTP 配置为空")
		return false
	}

	if nil == a.cfg.Http.Basicauth {
		a.logger.Error("认证失败：基本认证配置为空")
		return false
	}

	// 如果未启用基本认证，则允许所有请求通过。
	if !a.cfg.Http.Basicauth.Enabled {
		a.logger.Debug("基本认证未启用，请求自动通过")
		return true
	}

	// 检查用户名和密码是否匹配任何一个凭证。
	var userFound bool
	var userCredential *config.Credential

	// 首先检查用户名是否存在。
	for _, credential := range a.cfg.Http.Basicauth.Credentials {
		if credential.Username == username {
			userFound = true
			userCredential = credential
			break
		}
	}

	// 如果用户名不存在，记录具体错误。
	if !userFound {
		a.logger.WithField("username", username).Warn("认证失败：用户名不存在")
		return false
	}

	// 用户名存在，检查密码是否正确。
	if userCredential.Password != password {
		a.logger.WithField("username", username).Warn("认证失败：密码不匹配")
		return false
	}

	// 认证成功。
	a.logger.WithField("username", username).Debug("认证成功")

	return true
}
