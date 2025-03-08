// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package biz 实现业务逻辑层，包含领域模型和业务规则。
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
	// 确保 smsBiz 实现了 SmsBiz 接口。
	_ SmsBiz = (*smsBiz)(nil)
)

type (
	// SmsInfo 定义短信信息结构。
	SmsInfo struct {
		// 短信唯一标识符。
		ID string `json:"id"`
		// 发送方手机号。
		From string `json:"from"`
		// 接收方手机号。
		To string `json:"to"`
		// 短信内容。
		Message string `json:"message"`
		// 创建时间。
		CreatedAt time.Time `json:"created_at"`
	}

	// smsBiz 实现了 SmsBiz 接口，提供短信业务逻辑。
	smsBiz struct {
		// 日志记录器。
		logger log.Logger
		// 应用配置。
		cfg *config.Config
		// 短信仓储接口。
		repo SmsRepo
	}
)

// NewSmsBiz 创建短信业务逻辑实例。
//
// 参数：
//   - logger log.Logger：日志记录器。
//   - cfg *config.Config：应用配置。
//   - repo SmsRepo：短信仓储接口。
//
// 返回值：
//   - SmsBiz：短信业务逻辑接口实例。
func NewSmsBiz(logger log.Logger, cfg *config.Config, repo SmsRepo) SmsBiz {
	return &smsBiz{
		logger: logger.WithField("ddd", "biz").WithField("module", "sms"),
		cfg:    cfg,
		repo:   repo,
	}
}

// SendSms 实现发送短信的业务逻辑。
// 将短信信息保存到仓储中。
//
// 参数：
//   - ctx context.Context：上下文。
//   - sms *SmsInfo：短信信息。
//
// 返回值：
//   - error：处理过程中可能发生的错误。
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

// List 实现查询短信列表的业务逻辑。
// 从仓储中获取所有短信记录。
//
// 参数：
//   - ctx context.Context：上下文。
//
// 返回值：
//   - []*SmsInfo：短信记录列表。
//   - error：处理过程中可能发生的错误。
func (s *smsBiz) List(ctx context.Context) ([]*SmsInfo, error) {
	return s.repo.List(ctx)
}
