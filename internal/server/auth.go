// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package server

import (
	"context"
)

// 用户名和密码的映射，实际项目中应该从数据库或配置文件中读取
var authUsers = map[string]string{
	"admin": "admin", // 仅用于演示，生产环境请使用更复杂的密码
}

// Authenticate 验证用户名和密码是否正确
func Authenticate(ctx context.Context, username, password string) bool {
	if storedPassword, exists := authUsers[username]; exists {
		return storedPassword == password
	}
	return false
}
