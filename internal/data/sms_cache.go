package data

import (
	"time"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/biz"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

type (
	SmsCache struct {
		logger log.Logger
		cfg    *config.Config

		c      map[string]*biz.SmsInfo
		ticker *time.Ticker
		done   chan bool
	}
)

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

func (s *SmsCache) Set(key string, value *biz.SmsInfo) {
	s.c[key] = value
}

func (s *SmsCache) Get(key string) *biz.SmsInfo {
	return s.c[key]
}

// 清理过期记录，删除创建时间超过 1 小时的记录。
func (s *SmsCache) cleanupExpiredRecords() {
	now := time.Now()
	// 指定秒数前的时间。
	expireTime := now.Add(-1 * time.Duration(s.cfg.Sms.Cache.Time) * time.Second)

	for key, info := range s.c {
		// 检查记录是否已经超过 1 小时。
		if info.CreatedAt.Before(expireTime) {
			delete(s.c, key)
			s.logger.
				WithField("id", info.ID).
				Debug("清理过期短信记录")
		}
	}
}

func (s *SmsCache) cleanupSmsCache() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
	close(s.done)
}
