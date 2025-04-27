package common

import "time"

// CANP协议常量
const (
	CANPHederLen       = 22              // CANP报文头部长度
	CANPTypePublish    = 0x01            // 发布报文类型
	CANPTypeQuery      = 0x02            // 查询报文类型
	DefaultPort        = 38000           // 天基边缘云默认监听端口
	DefaultCarrierIP   = "127.0.0.1"     // 承载网默认监听端口
	DefaultCarrierPort = 58000           // 承载网默认监听端口
	SendPublishPeriod  = 2 * time.Second // 默认发送周期
	LocalHost          = "127.0.0.1"
	NodeID             = 1024
	NodeAddress        = "2001:db8::1"
	ReadTimeout        = 1000 * time.Microsecond
)
