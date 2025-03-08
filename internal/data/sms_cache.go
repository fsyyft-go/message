// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package data 实现数据访问层，负责数据的存储和检索。
package data

import (
	"sort"
	"time"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/message/internal/biz"
	"github.com/fsyyft-go/message/internal/config"
)

type (
	// SmsCache 实现短信缓存功能，提供短信数据的内存存储和检索。
	SmsCache struct {
		// 日志记录器。
		logger log.Logger
		// 应用配置。
		cfg *config.Config

		// 短信数据缓存，键为短信 ID，值为短信信息。
		c map[string]*biz.SmsInfo
		// 定时器，用于定期清理过期记录。
		ticker *time.Ticker
		// 停止信号通道，用于优雅关闭定时器。
		done chan bool
	}
)

// NewSmsCache 创建短信缓存实例。
//
// 参数：
//   - logger log.Logger：日志记录器。
//   - cfg *config.Config：应用配置。
//
// 返回值：
//   - *SmsCache：短信缓存实例。
//   - func()：清理函数，用于资源释放。
//   - error：初始化过程中可能发生的错误。
func NewSmsCache(logger log.Logger, cfg *config.Config) (*SmsCache, func(), error) {
	ct := &SmsCache{
		cfg:    cfg,
		c:      make(map[string]*biz.SmsInfo),
		logger: logger.WithField("ddd", "data").WithField("module", "sms").WithField("do", "sms cache"),
		done:   make(chan bool),
	}

	ct.logger.Info("初始化短信缓存")

	// 启动定时任务，定时清理一次过期记录。
	ct.ticker = time.NewTicker(time.Second * time.Duration(ct.cfg.Sms.Cache.Interval))
	go func() {
		ct.logger.WithField("interval", ct.cfg.Sms.Cache.Interval).Info("启动清理过期记录定时任务")
		for {
			select {
			case <-ct.ticker.C:
				ct.cleanupExpiredRecords()
			case <-ct.done:
				ct.logger.Info("停止清理过期记录")
				return
			}
		}
	}()

	return ct, func() { ct.cleanupSmsCache() }, nil
}

// Set 设置短信缓存，将短信信息存储到缓存中。
//
// 参数：
//   - key string：短信 ID。
//   - value *biz.SmsInfo：短信信息。
func (s *SmsCache) Set(key string, value *biz.SmsInfo) {
	s.c[key] = value
}

// List 获取所有短信记录，并按创建时间倒序排序。
//
// 返回值：
//   - []*biz.SmsInfo：短信记录列表。
func (s *SmsCache) List() []*biz.SmsInfo {
	infos := make([]*biz.SmsInfo, 0, len(s.c))
	for _, info := range s.c {
		infos = append(infos, info)
	}
	// 按时间倒序排序。
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].CreatedAt.After(infos[j].CreatedAt)
	})
	return infos
}

// cleanupExpiredRecords 清理过期记录，删除创建时间超过配置的过期时间的记录。
func (s *SmsCache) cleanupExpiredRecords() {
	now := time.Now()
	// 指定秒数前的时间。
	expireTime := now.Add(-1 * time.Duration(s.cfg.Sms.Cache.Time) * time.Second)

	for key, info := range s.c {
		// 检查记录是否已经超过配置的过期时间。
		if info.CreatedAt.Before(expireTime) {
			delete(s.c, key)
			s.logger.
				WithField("id", info.ID).
				Debug("清理过期短信记录")
		}
	}
}

// cleanupSmsCache 清理短信缓存资源，停止定时器并关闭通道。
func (s *SmsCache) cleanupSmsCache() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
	close(s.done)
}
