// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package biz

import (
	"context"
	"time"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

type (
	SmsInfo struct {
		ID        string    `json:"id"`
		From      string    `json:"from"`
		To        string    `json:"to"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"created_at"`
	}

	SmsRepo interface {
		Save(ctx context.Context, sms *SmsInfo) error
	}
)

type SmsBiz struct {
	logger log.Logger
	cfg    *config.Config
	repo   SmsRepo
}

func NewSmsBiz(logger log.Logger, cfg *config.Config, repo SmsRepo) *SmsBiz {
	return &SmsBiz{
		logger: logger.WithField("ddd", "biz").WithField("module", "sms"),
		cfg:    cfg,
		repo:   repo,
	}
}

func (s *SmsBiz) SendSms(ctx context.Context, sms *SmsInfo) error {
	l := s.logger

	if nil != sms {
		l = l.WithField("from", sms.From).WithField("to", sms.To).WithField("message", sms.Message).WithField("id", sms.ID)
		l.Debug("")
	} else {
		l.Warn("请求参数为空")
	}

	if err := s.repo.Save(ctx, sms); nil != err {
		l.WithField("error", err).Error("保存短信失败")
		return err
	}

	return nil
}
