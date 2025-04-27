package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"canp.server/canp"
	"canp.server/canp_cfg"
	"canp.server/carrier"
	"canp.server/common"
	"canp.server/edge"
	"canp.server/timer"
)

func main() {
	// 读取配置文件
	cfg, err := canp_cfg.ReadConfigCurrentDir("canp_cfg.yaml")
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	// 资源查询通道
	queryChan := make(chan *canp.QueryRequest, 100)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go edge.StartRcvSendPublishPacket(common.LocalHost, cfg.EdgeInfo.ServerPort, queryChan, ctx)
	go timer.StartTimer(common.SendPublishPeriod, queryChan, ctx, cfg.CarrierInfo.ServerIP, cfg.CarrierInfo.ServerPort)

	// 发送查询报文和接收发布报文，主要使用本地carrier进行测试
	go carrier.StartRcvSendPacket(cfg.CarrierInfo.ServerIP, cfg.CarrierInfo.ServerPort, ctx, common.LocalHost, cfg.EdgeInfo.ServerPort)

	// 设置信号监听
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	// 等待信号
	<-sigCh
	fmt.Println("\n接收到中断信号, 开始清理...")

	// 取消上下文，通知所有goroutine开始清理
	cancel()

	// 等待所有工作完成
	time.Sleep(10 * time.Second)
	fmt.Println("CANP server end")
}
