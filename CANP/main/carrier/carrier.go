// carrier/carrier.go
package carrier

import (
	"context"
	"fmt"
	"net"
	"time"

	"canp.server/canp"
	"canp.server/common"
)

// 发送查询请求协程
func startQuerySender(conn *net.UDPConn, dstIP net.IP, dstPort int) {
	// 解析目标地址
	target, err := net.ResolveUDPAddr("udp", dstIP.String()+fmt.Sprintf(":%d", dstPort))
	if err != nil {
		fmt.Printf("carrier地址解析失败: %v\n", err)
	}

	for i := 0; i < 1; i++ {
		packet, err := canp.CreateQueryRequest(
			common.NodeID,                   // 查询节点ID
			net.ParseIP(common.NodeAddress), // 查询节点地址
			[]uint8{0x01, 0x02},             // 查询类型
		)
		if err != nil {
			fmt.Printf("carrier构建查询请求失败: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		// 通过同一个连接发送
		if _, err := conn.WriteToUDP(packet, target); err != nil {
			fmt.Printf("carrier发送失败: %v\n", err)
		} else {
			fmt.Printf("carrier[发送] 到 %s:%d 成功\n", target.IP, target.Port)
		}

		fmt.Println("carrier已发送状态查询报文")
		time.Sleep(5 * time.Second) // 每5秒发送一次查询
	}
}

// 接收状态发布报文协程
func startReceiver(conn *net.UDPConn, ctx context.Context) {

	buffer := make([]byte, 65535)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("carrier接收到结束通知, 停止接收")
			return
		default:
			conn.SetReadDeadline(time.Now().Add(common.ReadTimeout)) // 1秒超时
			n, addr, err := conn.ReadFromUDP(buffer)
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 处理超时逻辑, 超时不做任何操作
				continue
			} else if err != nil {
				fmt.Printf("carrier接收报文失败: %v\n", err)
				continue
			}

			handleReceivedPublishPacket(buffer[:n], addr)
		}
	}
}

// 处理接收到的状态发布报文
func handleReceivedPublishPacket(data []byte, addr *net.UDPAddr) {
	fmt.Printf("carrier接收到数据包长度: %d, 值: %v\n", len(data), data)
	header, tlvs, err := canp.ParsePublishPacket(data)
	if err != nil {
		fmt.Printf("carrier解析报文失败: %v\n", err)
		return
	}

	fmt.Printf("=== carrier收到状态发布报文 ===\n")
	fmt.Printf("源IP: %s\n", addr.IP)
	fmt.Printf("源端口: %d\n", addr.Port)
	fmt.Printf("节点ID: %d\n", header.NodeID)
	fmt.Printf("节点地址: %s\n", header.NodeAddr)
	fmt.Printf("TLV数据:\n")

	for tlvType, tlvValue := range tlvs {
		fmt.Printf("  类型: 0x%02x 长度: %d 值: %v\n", tlvType, tlvValue.Length, tlvValue.Value)
	}
}

func StartRcvSendPacket(listenIp string, listenPort int, ctx context.Context, sendIp string, sendPort int) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(listenIp),
		Port: listenPort,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Printf("carrier接收服务已启动，监听地址: %s:%d\n", listenIp, listenPort)
	// 启动接收协程
	go startReceiver(conn, ctx)

	// startQuerySender(conn, net.ParseIP(sendIp), sendPort)
	fmt.Println("carrier发送服务已启动")

	// 用于提前通知停止
	<-ctx.Done()
	fmt.Println("carrier服务已停止")
}
