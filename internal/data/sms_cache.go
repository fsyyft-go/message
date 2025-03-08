package data

import (
	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/biz"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

type (
	SmsCache struct {
		c map[string]*biz.SmsInfo
	}
)

func NewSmsCache(logger log.Logger, cfg *config.Config) (*SmsCache, func(), error) {
	_ = cfg

	ct := &SmsCache{c: make(map[string]*biz.SmsInfo)}

	logger.WithField("ddd", "data").WithField("module", "sms").Info("初始化短信缓存")

	return ct, cleanupSmsCache, nil
}

func (s *SmsCache) Set(key string, value *biz.SmsInfo) {
	s.c[key] = value
}

func (s *SmsCache) Get(key string) *biz.SmsInfo {
	return s.c[key]
}

func cleanupSmsCache() {
}
