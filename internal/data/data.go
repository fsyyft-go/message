// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package data 实现数据访问层，负责数据的存储和检索。
package data

import (
	"github.com/google/wire"
)

var (
	// ProviderSet 是数据访问层的依赖注入提供者集合。
	// 包含短信缓存和短信仓储的创建函数。
	ProviderSet = wire.NewSet(
		NewSmsCache,
		NewSmsRepo,
	)
)
