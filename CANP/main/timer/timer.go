package timer

import (
	"context"
	"fmt"
	"net"
	"time"

	"canp.server/canp"
	"canp.server/common"
)

// 周期性任务调度器
type PeriodicScheduler struct {
	interval   time.Duration
	task       func() (*canp.QueryRequest, error)
	resultChan chan<- *canp.QueryRequest
	stopCtx    context.Context
}

// 创建新的周期性任务调度器
func NewPeriodicScheduler(
	interval time.Duration,
	task func() (*canp.QueryRequest, error),
	queryChan chan<- *canp.QueryRequest,
	stopCtx context.Context,
) *PeriodicScheduler {
	return &PeriodicScheduler{
		interval:   interval,
		task:       task,
		resultChan: queryChan,
		stopCtx:    stopCtx,
	}
}

// 启动定时器
func (s *PeriodicScheduler) Start() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	fmt.Println("定时器已启动，间隔时间:", s.interval)
	for {
		select {
		case <-ticker.C:
			query, err := s.task()
			if err != nil {
				fmt.Println("Error executing task:", err)
			} else {
				s.resultChan <- query
			}
		case <-s.stopCtx.Done():
			fmt.Println("定时器接收到停止信号，停止调度")
			return
		}
	}
}

func StartTimer(
	interval time.Duration,
	queryChan chan<- *canp.QueryRequest,
	ctx context.Context,
	carrier_ip string,
	carrier_port int,
) {

	task := func() (*canp.QueryRequest, error) {
		// 这里可以添加你需要执行的任务逻辑
		// 例如，发送查询请求
		return &canp.QueryRequest{
			Header: &canp.CANPHeader{
				Type:     common.CANPTypePublish,
				NodeID:   common.NodeID,
				NodeAddr: net.IPv6loopback,
			},
			SrcIP:        net.ParseIP(carrier_ip),
			SrcPort:      uint16(carrier_port),
			PeriodicSend: true,
			TLVTypes:     []uint8{}, // 这里可以添加你需要查询的TLV类型
		}, nil
	}

	scheduler := NewPeriodicScheduler(interval, task, queryChan, ctx)
	go scheduler.Start()
}
