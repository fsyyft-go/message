// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/api/sms"
	"github.com/fsyyft-go/sms-bridge/internal/biz"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

var (
	_ sms.SmsServiceHTTPServer = (*SmsService)(nil)
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
		ID:        uuid.New().String(),
		From:      req.From,
		To:        req.To,
		Message:   req.Message,
		CreatedAt: time.Now(),
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

func (s *SmsService) ListSms(ctx context.Context, req *sms.ListSmsRequest) (*sms.ListSmsResponse, error) {
	infos, err := s.biz.List(ctx)
	if nil != err {
		s.logger.WithField("error", err).Error("获取短信列表失败")
		return nil, err
	}

	smsInfos := make([]*sms.SmsInfo, 0, len(infos))
	for _, info := range infos {
		smsInfos = append(smsInfos, &sms.SmsInfo{
			From:      info.From,
			To:        info.To,
			Message:   info.Message,
			MessageId: info.ID,
			CreatedAt: info.CreatedAt.Format(time.RFC3339),
		})
	}

	return &sms.ListSmsResponse{
		Sms: smsInfos,
	}, nil
}
