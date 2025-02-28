// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/api/sms"
	"github.com/fsyyft-go/sms-bridge/internal/biz"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

type (
	SmsService struct {
		logger log.Logger
		cfg    *config.Config

		biz *biz.SmsBiz
	}
)

func NewSmsService(logger log.Logger, cfg *config.Config, biz *biz.SmsBiz) *SmsService {
	return &SmsService{
		logger: logger.WithField("ddd", "service").WithField("module", "sms"),
		cfg:    cfg,
		biz:    biz,
	}
}

func (s *SmsService) SendSms(ctx context.Context, req *sms.SendSmsRequest) (*sms.SendSmsResponse, error) {
	l := s.logger

	if nil != req {
		l = l.WithField("from", req.From).WithField("to", req.To).WithField("message", req.Message)
		l.Debug("")
	} else {
		l.Warn("请求参数为空")
	}

	smsInfo := &biz.SmsInfo{
		ID:      uuid.New().String(),
		From:    req.From,
		To:      req.To,
		Message: req.Message,
	}

	err := s.biz.SendSms(ctx, smsInfo)

	if nil != err {
		l.WithField("id", smsInfo.ID).WithField("error", err).Error("发送短信失败")
		return nil, err
	}

	return &sms.SendSmsResponse{
		MessageId: smsInfo.ID,
	}, nil
}
