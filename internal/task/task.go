// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package task

import (
	"context"

	"github.com/google/wire"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/message/internal/config"
	"github.com/fsyyft-go/message/pkg/runtime"
)

var (
	ProviderSet = wire.NewSet(
		NewTask,
	)

	_ runtime.Runner = (*task)(nil)
	_ Task           = (*task)(nil)
)

type (
	// Task 定义了任务接口。
	// 它继承了 runtime.Runner 接口，并添加了 TaskCount 方法。
	Task interface {
		runtime.Runner
		// TaskCount 返回任务数量。
		TaskCount() int
	}

	task struct {
		logger log.Logger
		cfg    *config.Config
	}
)

func NewTask(logger log.Logger, cfg *config.Config) Task {
	return &task{
		logger: logger,
		cfg:    cfg,
	}
}

func (t *task) Start(ctx context.Context) error {
	return nil
}

func (t *task) Stop(ctx context.Context) error {
	return nil
}

func (t *task) TaskCount() int {
	return 0
}
