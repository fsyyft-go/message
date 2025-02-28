package sms_bridge

import (
	"flag"
	"fmt"

	"github.com/fsyyft-go/kit/log"
	"github.com/fsyyft-go/sms-bridge/internal/config"
)

func Run() {
	var configPath string

	flag.StringVar(&configPath, "config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("加载配置文件失败：%v", err)
	}

	// 初始化日志。
	if err := log.InitLogger(
		log.WithLogType(log.LogType(cfg.Log.Type)),
		log.WithOutput(cfg.Log.Output),
	); err != nil {
		panic(err)
	}

	if webServer, cleanup, err := wireServer(log.GetLogger(), cfg); err != nil {
		cleanup()
	} else {
		_ = webServer.Start()
	}
}
