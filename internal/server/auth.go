// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package server 实现服务器层，负责 HTTP 服务器的配置和启动。
package server

import (
	"context"
)

// 用户名和密码的映射，实际项目中应该从数据库或配置文件中读取。
var authUsers = map[string]string{
	"admin": "admin", // 仅用于演示，生产环境请使用更复杂的密码。
}

// Authenticate 验证用户名和密码是否正确。
// 用于基本认证中间件的验证函数。
//
// 参数：
//   - ctx context.Context：上下文。
//   - username string：用户名。
//   - password string：密码。
//
// 返回值：
//   - bool：如果认证成功，返回 true；否则返回 false。
func Authenticate(ctx context.Context, username, password string) bool {
	if storedPassword, exists := authUsers[username]; exists {
		return storedPassword == password
	}
	return false
}
