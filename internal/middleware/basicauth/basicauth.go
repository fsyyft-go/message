// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

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
	Username string
	Password string
}

// Option 是 basicauth 中间件的配置选项。
type Option func(*options)

// options 包含中间件配置选项。
type options struct {
	validator CredentialValidator
	realm     string
}

// CredentialValidator 是一个函数类型，用于验证认证信息。
type CredentialValidator func(ctx context.Context, username, password string) bool

// WithValidator 配置自定义的认证信息验证器。
func WithValidator(validator CredentialValidator) Option {
	return func(o *options) {
		o.validator = validator
	}
}

// WithRealm 配置认证域名，显示在浏览器认证对话框中。
func WithRealm(realm string) Option {
	return func(o *options) {
		o.realm = realm
	}
}

// Server 创建一个基本认证中间件，用于服务端验证请求。
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
					// 设置 WWW-Authenticate 头，触发浏览器的认证对话框
					tr.ReplyHeader().Set("WWW-Authenticate", `Basic realm="`+o.realm+`"`)
					return nil, ErrInvalidBasicAuth
				}

				// 解析认证头
				cred, err := parseBasicAuth(auths)
				if err != nil {
					// 设置 WWW-Authenticate 头，触发浏览器的认证对话框
					tr.ReplyHeader().Set("WWW-Authenticate", `Basic realm="`+o.realm+`"`)
					return nil, ErrInvalidBasicAuth
				}

				// 验证用户名和密码
				if !o.validator(ctx, cred.Username, cred.Password) {
					// 设置 WWW-Authenticate 头，触发浏览器的认证对话框
					tr.ReplyHeader().Set("WWW-Authenticate", `Basic realm="`+o.realm+`"`)
					return nil, ErrInvalidBasicAuth
				}
			}
			return handler(ctx, req)
		}
	}
}

// parseBasicAuth 解析 HTTP Basic Auth 头部值，返回认证信息。
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
