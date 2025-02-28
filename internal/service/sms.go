// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package service

import (
	"context"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/api/sms"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

type (
	SmsService struct {
		logger log.Logger
		cfg    *config.Config
	}
)

func NewSmsService(logger log.Logger, cfg *config.Config) *SmsService {
	return &SmsService{
		logger: logger,
		cfg:    cfg,
	}
}

func (s *SmsService) SendSms(ctx context.Context, req *sms.SendSmsRequest) (*sms.SendSmsResponse, error) {
	s.logger.Debug("SendSms", "from", req.From, "to", req.To, "message", req.Message)
	return &sms.SendSmsResponse{}, nil
}
