package server

import (
	"fmt"
	"net"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/api/sms"
	"github.com/fsyyft-go/sms-bridge/internal/config"
	"github.com/fsyyft-go/sms-bridge/internal/service"
	bridge_kratos_http "github.com/fsyyft-go/sms-bridge/pkg/kratos/transport/http"
)

type (
	WebServer struct {
		logger log.Logger
		cfg    *config.Config
		server *http.Server
	}
)

func NewWebServer(logger log.Logger, cfg *config.Config, smsService *service.SmsService) (*WebServer, func(), error) {
	var err error
	webServer := &WebServer{
		logger: logger,
		cfg:    cfg,
	}

	server := http.NewServer(http.Middleware(
		recovery.Recovery(),
		validate.Validator(),
	))
	sms.RegisterSmsServiceHTTPServer(server, smsService)

	fmt.Println(bridge_kratos_http.GetPaths(server))

	webServer.server = server

	var cleanup = func() {}

	return webServer, cleanup, err
}

func (s *WebServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Http.Port))
	if err != nil {
		return err
	}

	return s.server.Serve(listener)
}
