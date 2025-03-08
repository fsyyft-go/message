// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package sms_bridge 实现了短信网桥服务的核心功能。
package sms_bridge

import (
	"flag"
	"fmt"

	"github.com/google/wire"

	"github.com/fsyyft-go/message/internal/config"
)

var (
	// ProviderSet 是一个 Wire 依赖注入提供者集合，用于注册可被注入的组件。
	// 当前仅包含日志记录器的创建函数。
	ProviderSet = wire.NewSet(
		NewLogger,
	)
)

// Run 函数是应用程序的入口点，负责启动短信网桥服务。
// 它完成以下任务：
// 1. 解析命令行参数获取配置文件路径。
// 2. 加载应用配置。
// 3. 通过 Wire 框架初始化所有依赖。
// 4. 启动 Web 服务器。
func Run() {
	// 定义配置文件路径变量，默认为"configs/config.yaml"。
	var configPath string

	// 注册命令行参数，用于指定配置文件路径。
	flag.StringVar(&configPath, "config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	// 从指定路径加载配置文件。
	cfg, err := config.LoadConfig(configPath)
	if nil != err {
		fmt.Printf("加载配置文件失败：%v", err)
		return
	}

	// 通过 Wire 框架生成的 wireServer 函数初始化服务。
	// 该函数会自动注入所有依赖项并返回配置好的 Web 服务器实例。
	if webServer, cleanup, err := wireServer(cfg); nil != err {
		fmt.Printf("初始化失败：%v", err)
		// 调用清理函数释放已分配的资源。
		cleanup()
	} else {
		// 启动 Web 服务器。
		_ = webServer.Start()
	}
}
