// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package config 提供应用程序配置的加载和管理功能。
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
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
	// 从 yaml 文件中读取配置。
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%v", err)
	}

	// 解析 yaml 配置到结构体。
	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败：%v", err)
	}

	return &cfg, nil
}
