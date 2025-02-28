package main

import (
	"flag"
	"fmt"

	"github.com/fsyyft-go/sms-bridge/internal/config"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("加载配置文件失败：%v", err)
	}
	fmt.Println(cfg)
}
