// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package config 提供应用程序配置的加载和管理功能。
package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	kit_kratos_config "github.com/fsyyft-go/kit/kratos/config"
)

// LoadConfig 从指定路径加载配置文件并解析为 Config 结构体。
//
// 参数：
//   - path string：配置文件的路径。
//
// 返回值：
//   - *Config：解析后的配置对象指针。
//   - error：加载或解析过程中可能发生的错误。
func LoadConfig(path string) (*Config, error) {
	// 创建配置管理器实例。
	c := config.New(
		// 设置配置源为文件源，指定配置文件路径。
		config.WithSource(
			file.NewSource(path),
		),
		// 设置自定义解码器，支持特殊格式处理（如 base64 解码）。
		config.WithDecoder(kit_kratos_config.NewDecoder().Decode),
	)
	// 加载配置，如果出错则触发 panic。
	if err := c.Load(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := c.Scan(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
