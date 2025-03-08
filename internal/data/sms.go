// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package data 实现数据访问层，负责数据的存储和检索。
package data

import (
	"context"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/message/internal/biz"
	"github.com/fsyyft-go/message/internal/config"
)

var (
	// 确保 Sms 实现了 biz.SmsRepo 接口。
	_ biz.SmsRepo = (*Sms)(nil)
)

type (
	// Sms 实现了 biz.SmsRepo 接口，提供短信数据的存储和检索功能。
	Sms struct {
		// 日志记录器。
		logger log.Logger
		// 应用配置。
		cfg *config.Config
		// 短信缓存。
		cache *SmsCache
	}
)

// NewSmsRepo 创建短信仓储实例。
//
// 参数：
//   - logger log.Logger：日志记录器。
//   - cfg *config.Config：应用配置。
//   - cache *SmsCache：短信缓存。
//
// 返回值：
//   - biz.SmsRepo：短信仓储接口实例。
func NewSmsRepo(logger log.Logger, cfg *config.Config, cache *SmsCache) biz.SmsRepo {
	return &Sms{
		logger: logger.WithField("ddd", "data").WithField("module", "sms"),
		cfg:    cfg,
		cache:  cache,
	}
}

// Save 实现保存短信记录的功能。
// 将短信信息保存到缓存中。
//
// 参数：
//   - ctx context.Context：上下文。
//   - sms *biz.SmsInfo：短信信息。
//
// 返回值：
//   - error：处理过程中可能发生的错误。
func (s *Sms) Save(ctx context.Context, sms *biz.SmsInfo) error {
	s.cache.Set(sms.ID, sms)

	return nil
}

// List 实现查询短信记录列表的功能。
// 从缓存中获取所有短信记录。
//
// 参数：
//   - ctx context.Context：上下文。
//
// 返回值：
//   - []*biz.SmsInfo：短信记录列表。
//   - error：处理过程中可能发生的错误。
func (s *Sms) List(ctx context.Context) ([]*biz.SmsInfo, error) {
	return s.cache.List(), nil
}
