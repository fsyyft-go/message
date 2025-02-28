// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package data

import (
	"context"
	"errors"

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
	}
)

func NewSmsRepo(logger log.Logger, cfg *config.Config) biz.SmsRepo {
	return &Sms{
		logger: logger.WithField("ddd", "data").WithField("module", "sms"),
		cfg:    cfg,
	}
}

func (s *Sms) Save(ctx context.Context, sms *biz.SmsInfo) error {
	return errors.New("not implemented")
}
