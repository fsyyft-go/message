// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package main 是应用程序的入口包。
package main

import (
	sms "github.com/fsyyft-go/sms-bridge/internal/app/sms-bridge"
)

// main 函数是应用程序的入口点，调用 sms 包中的 Run 函数启动短信网桥服务。
func main() {
	sms.Run()
}
