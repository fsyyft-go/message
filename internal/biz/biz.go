// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package biz 实现业务逻辑层，包含领域模型和业务规则。
package biz

import (
	"github.com/google/wire"
)

var (
	// ProviderSet 是业务逻辑层的依赖注入提供者集合。
	// 包含短信业务逻辑的创建函数。
	ProviderSet = wire.NewSet(
		NewSmsBiz,
	)
)
