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
	return &SmsService{}
}

func (s *SmsService) SendSms(ctx context.Context, req *sms.SendSmsRequest) (*sms.SendSmsResponse, error) {
	return nil, nil
}
