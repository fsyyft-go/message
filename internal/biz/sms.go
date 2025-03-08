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

// 仓储接口定义开始。

type (
	// SmsRepo 定义了短信存储的仓储接口。
	// 该接口提供了保存和查询短信记录的基本功能。
	// 具体实现可以是内存缓存、数据库等多种形式。
	SmsRepo interface {
		// Save 保存短信记录。
		// 参数:
		//   ctx: 上下文。
		//   sms: 短信记录。
		// 返回:
		//   错误信息。
		Save(ctx context.Context, sms *SmsInfo) error

		// List 查询短信记录。
		// 参数:
		//   ctx: 上下文。
		// 返回:
		//   短信记录列表。
		List(ctx context.Context) ([]*SmsInfo, error)
	}
)

// 仓储接口定义结束。

// 领域接口定义开始。

type (
	// SmsBiz 定义了短信业务的领域接口。
	// 该接口提供了发送短信和查询短信记录的基本功能。
	// 具体实现可以是内存缓存、数据库等多种形式。
	SmsBiz interface {
		// SendSms 发送短信。
		// 参数:
		//   ctx: 上下文。
		//   sms: 短信记录。
		// 返回:
		//   错误信息。
		SendSms(ctx context.Context, sms *SmsInfo) error

		// List 查询短信记录。
		// 参数:
		//   ctx: 上下文。
		// 返回:
		//   短信记录列表。
		List(ctx context.Context) ([]*SmsInfo, error)
	}
)

// 领域接口定义结束。

var (
	_ SmsBiz = (*smsBiz)(nil)
)

type (
	SmsInfo struct {
		ID        string    `json:"id"`
		From      string    `json:"from"`
		To        string    `json:"to"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"created_at"`
	}

	smsBiz struct {
		logger log.Logger
		cfg    *config.Config
		repo   SmsRepo
	}
)

func NewSmsBiz(logger log.Logger, cfg *config.Config, repo SmsRepo) SmsBiz {
	return &smsBiz{
		logger: logger.WithField("ddd", "biz").WithField("module", "sms"),
		cfg:    cfg,
		repo:   repo,
	}
}

func (s *smsBiz) SendSms(ctx context.Context, sms *SmsInfo) error {
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

func (s *smsBiz) List(ctx context.Context) ([]*SmsInfo, error) {
	return s.repo.List(ctx)
}
