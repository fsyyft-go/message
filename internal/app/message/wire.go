// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// 构建标签说明。
// 该构建标签确保此存根代码不会包含在最终构建中。
//go:build wireinject
// +build wireinject

// package sms_bridge 实现了短信网桥服务的核心功能。
package sms_bridge

import (
	"github.com/google/wire"

	"github.com/fsyyft-go/message/internal/biz"
	"github.com/fsyyft-go/message/internal/config"
	"github.com/fsyyft-go/message/internal/data"
	"github.com/fsyyft-go/message/internal/server"
	"github.com/fsyyft-go/message/internal/service"
)

// wireServer 函数用于构建和初始化 WebServer 实例。
// 该函数使用 Google Wire 进行依赖注入，自动组装应用程序组件。
//
// 参数：
//   - cfg *config.Config：应用程序配置对象。
//
// 返回值：
//   - server.WebServer：初始化后的 Web 服务器实例。
//   - func()：清理函数，用于资源释放。
//   - error：初始化过程中可能发生的错误。
func wireServer(cfg *config.Config) (server.WebServer, func(), error) {
	// wire.Build 函数用于声明依赖关系图，将所有组件连接在一起。
	// panic 调用会在编译时被 wire 工具替换为实际的依赖注入代码。
	panic(wire.Build(
		// 引入当前包中定义的提供者集合。
		ProviderSet,
		// 引入数据层的提供者集合。
		data.ProviderSet,
		// 引入业务逻辑层的提供者集合。
		biz.ProviderSet,
		// 引入服务层的提供者集合。
		service.ProviderSet,
		// 引入服务器层的提供者集合。
		server.ProviderSet,
	))
}
