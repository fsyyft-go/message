//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package sms_bridge

import (
	"github.com/google/wire"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/config"
	"github.com/fsyyft-go/sms-bridge/internal/server"
	"github.com/fsyyft-go/sms-bridge/internal/service"
)

func wireServer(logger log.Logger, cfg *config.Config) (*server.WebServer, func(), error) {
	panic(wire.Build(
		service.ProviderSet,
		server.ProviderSet,
	))
}
