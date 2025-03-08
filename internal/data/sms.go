// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package data

import (
	"context"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/biz"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

var (
	_ biz.SmsRepo = (*Sms)(nil)
)

type (
	Sms struct {
		logger log.Logger
		cfg    *config.Config

		cache *SmsCache
	}
)

func NewSmsRepo(logger log.Logger, cfg *config.Config, cache *SmsCache) biz.SmsRepo {
	return &Sms{
		logger: logger.WithField("ddd", "data").WithField("module", "sms"),
		cfg:    cfg,
		cache:  cache,
	}
}

func (s *Sms) Save(ctx context.Context, sms *biz.SmsInfo) error {
	s.cache.Set(sms.ID, sms)

	return nil
}
