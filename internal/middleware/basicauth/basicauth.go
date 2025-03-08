// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package basicauth 实现 HTTP 基本认证中间件，用于保护 API 接口。
package basicauth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// ErrInvalidBasicAuth 是无效的基本认证错误。
var ErrInvalidBasicAuth = errors.Unauthorized("UNAUTHORIZED", "Invalid basic authentication")

// Credential 保存用户认证信息。
type Credential struct {
	// 用户名。
	Username string
	// 密码。
	Password string
}

// Option 是 basicauth 中间件的配置选项。
type Option func(*options)

// options 包含中间件配置选项。
type options struct {
	// 认证信息验证器。
	validator CredentialValidator
	// 认证域，显示在浏览器认证对话框中。
	realm string
}

// CredentialValidator 是一个函数类型，用于验证认证信息。
//
// 参数：
//   - ctx context.Context：上下文。
//   - username string：用户名。
//   - password string：密码。
//
// 返回值：
//   - bool：如果认证成功，返回 true；否则返回 false。
type CredentialValidator func(ctx context.Context, username, password string) bool

// WithValidator 配置自定义的认证信息验证器。
//
// 参数：
//   - validator CredentialValidator：认证信息验证器。
//
// 返回值：
//   - Option：中间件配置选项。
func WithValidator(validator CredentialValidator) Option {
	return func(o *options) {
		o.validator = validator
	}
}

// WithRealm 配置认证域名，显示在浏览器认证对话框中。
//
// 参数：
//   - realm string：认证域名。
//
// 返回值：
//   - Option：中间件配置选项。
func WithRealm(realm string) Option {
	return func(o *options) {
		o.realm = realm
	}
}

// Server 创建一个基本认证中间件，用于服务端验证请求。
//
// 参数：
//   - opts ...Option：中间件配置选项。
//
// 返回值：
//   - middleware.Middleware：基本认证中间件。
func Server(opts ...Option) middleware.Middleware {
	o := &options{
		validator: func(ctx context.Context, username, password string) bool {
			// 默认认证器，返回 false，建议覆盖此函数
			return false
		},
		realm: "Restricted",
	}
	for _, opt := range opts {
		opt(o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				auths := tr.RequestHeader().Get("Authorization")
				if auths == "" {
					// 设置 WWW-Authenticate 头，触发浏览器的认证对话框。
					tr.ReplyHeader().Set("WWW-Authenticate", `Basic realm="`+o.realm+`"`)
					return nil, ErrInvalidBasicAuth
				}

				// 解析认证头。
				cred, err := parseBasicAuth(auths)
				if err != nil {
					// 设置 WWW-Authenticate 头，触发浏览器的认证对话框。
					tr.ReplyHeader().Set("WWW-Authenticate", `Basic realm="`+o.realm+`"`)
					return nil, ErrInvalidBasicAuth
				}

				// 验证用户名和密码。
				if !o.validator(ctx, cred.Username, cred.Password) {
					// 设置 WWW-Authenticate 头，触发浏览器的认证对话框。
					tr.ReplyHeader().Set("WWW-Authenticate", `Basic realm="`+o.realm+`"`)
					return nil, ErrInvalidBasicAuth
				}
			}
			return handler(ctx, req)
		}
	}
}

// parseBasicAuth 解析 HTTP Basic Auth 头部值，返回认证信息。
//
// 参数：
//   - auth string：Authorization 头部值。
//
// 返回值：
//   - *Credential：解析后的认证信息。
//   - error：解析过程中可能发生的错误。
func parseBasicAuth(auth string) (*Credential, error) {
	if !strings.HasPrefix(auth, "Basic ") {
		return nil, ErrInvalidBasicAuth
	}

	payload := strings.TrimPrefix(auth, "Basic ")
	decodedBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	decodedString := string(decodedBytes)
	parts := strings.SplitN(decodedString, ":", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidBasicAuth
	}

	return &Credential{
		Username: parts[0],
		Password: parts[1],
	}, nil
}
