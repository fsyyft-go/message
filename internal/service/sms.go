// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package service 实现服务层，负责处理 API 请求并调用业务逻辑。
package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/message/api/sms"
	"github.com/fsyyft-go/message/internal/biz"
	"github.com/fsyyft-go/message/internal/config"
)

var (
	// 确保 smsService 实现了 sms.SmsServiceHTTPServer 接口。
	_ sms.SmsServiceHTTPServer = (*smsService)(nil)
)

type (
	// smsService 实现了 sms.SmsServiceHTTPServer 接口，提供短信服务的 HTTP 接口。
	smsService struct {
		// 日志记录器。
		logger log.Logger
		// 应用配置。
		cfg *config.Config
		// 短信业务逻辑。
		biz biz.SmsBiz
	}
)

// NewSmsService 创建短信服务实例。
//
// 参数：
//   - logger log.Logger：日志记录器。
//   - cfg *config.Config：应用配置。
//   - biz biz.SmsBiz：短信业务逻辑。
//
// 返回值：
//   - sms.SmsServiceHTTPServer：短信服务 HTTP 接口实例。
func NewSmsService(logger log.Logger, cfg *config.Config, biz biz.SmsBiz) sms.SmsServiceHTTPServer {
	return &smsService{
		logger: logger.WithField("ddd", "service").WithField("module", "sms"),
		cfg:    cfg,
		biz:    biz,
	}
}

// SendSms 实现发送短信的 HTTP 接口。
// 接收客户端发送的短信请求，生成唯一标识符，并调用业务逻辑发送短信。
//
// 参数：
//   - ctx context.Context：上下文。
//   - req *sms.SendSmsRequest：发送短信请求。
//
// 返回值：
//   - *sms.SendSmsResponse：发送短信响应，包含短信唯一标识符。
//   - error：处理过程中可能发生的错误。
func (s *smsService) SendSms(ctx context.Context, req *sms.SendSmsRequest) (*sms.SendSmsResponse, error) {
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

// ListSms 实现查询短信列表的 HTTP 接口。
// 调用业务逻辑获取短信列表，并转换为 API 响应格式。
//
// 参数：
//   - ctx context.Context：上下文。
//   - req *sms.ListSmsRequest：查询短信列表请求。
//
// 返回值：
//   - *sms.ListSmsResponse：查询短信列表响应，包含短信列表。
//   - error：处理过程中可能发生的错误。
func (s *smsService) ListSms(ctx context.Context, req *sms.ListSmsRequest) (*sms.ListSmsResponse, error) {
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
