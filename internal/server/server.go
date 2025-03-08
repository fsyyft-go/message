// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package server 实现服务器层，负责 HTTP 服务器的配置和启动。
package server

import (
	"github.com/google/wire"
)

var (
	// ProviderSet 是服务器层的依赖注入提供者集合。
	// 包含 Web 服务器和认证器的创建函数。
	ProviderSet = wire.NewSet(
		NewAuthenticator,
		NewWebServer,
	)
)
