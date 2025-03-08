// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package sms_bridge 实现了短信网桥服务的核心功能。
package sms_bridge

import (
	"fmt"
	"os"
	"sync"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

var (
	// logger 是全局共享的日志记录器实例。
	logger log.Logger
	// loggerLocker 是用于保护 logger 变量的读写锁，确保并发安全。
	loggerLocker sync.RWMutex = sync.RWMutex{}
)

// NewLogger 创建并初始化一个日志记录器实例。
// 该函数使用单例模式确保只创建一个全局日志记录器。
//
// 参数：
//   - cfg *config.Config：应用程序配置对象，包含日志相关设置。
//
// 返回值：
//   - log.Logger：初始化后的日志记录器实例。
//   - func()：清理函数，用于在初始化失败时进行资源释放。
//   - error：初始化过程中可能发生的错误。
func NewLogger(cfg *config.Config) (log.Logger, func(), error) {
	var err error

	// 检查日志记录器是否已经初始化
	if nil == logger {
		// 加锁以防止并发初始化
		loggerLocker.Lock()
		defer loggerLocker.Unlock()

		// 双重检查锁定模式，确保日志记录器仅初始化一次。
		if nil == logger {
			// 使用配置创建新的日志记录器。
			if l, errNew := log.NewLogger(
				log.WithLogType(log.LogType(cfg.Log.Type)),
				log.WithOutput(cfg.Log.Output),
			); err != nil {
				err = errNew
			} else {
				// 设置日志级别。
				if level, err := log.ParseLevel(cfg.Log.Level); err != nil {
					l.WithField("error", err).Error("解析日志级别失败")
				} else {
					l.SetLevel(level)
					// 添加进程 ID 字段方便调试
					l = l.WithField("pid", os.Getpid())

					l.WithField("log_level", level).Info("设置日志级别")
				}

				// 更新全局日志记录器。
				logger = l
			}
		}
	}

	return logger, cleanupLogger, err
}

// cleanupLogger 清理日志记录器资源并记录初始化失败信息。
// 该函数在初始化出错时由 Wire 框架调用。
func cleanupLogger() {
	if nil != logger {
		// 如果日志记录器已初始化，使用它记录警告信息。
		logger.Warn("初始化失败")
		logger = nil
	} else {
		// 日志记录器未初始化，使用标准输出
		fmt.Println("初始化失败")
	}
}
