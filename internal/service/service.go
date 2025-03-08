// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package service 实现服务层，负责处理 API 请求并调用业务逻辑。
package service

import (
	"github.com/google/wire"
)

var (
	// ProviderSet 是服务层的依赖注入提供者集合。
	// 包含短信服务的创建函数。
	ProviderSet = wire.NewSet(
		NewSmsService,
	)
)
