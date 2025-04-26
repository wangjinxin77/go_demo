// carrier/carrier.go
package carrier

import (
	"fmt"
	"net"
	"time"

	"canp.server/canp"
)

// 发送查询请求协程
func StartQuerySender(dstIP string, dstPort uint16) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(dstIP),
		Port: int(dstPort),
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		packet, err := canp.CreateQueryRequest(
			net.ParseIP("192.168.0.2"), // 承载网源IP
			uint16(48000),              // 承载网源端口
			1024,                       // 查询节点ID
			net.ParseIP("172.10.0.1"),  // 查询节点地址
			[]uint8{0x01, 0x02},        // 查询类型
		)
		if err != nil {
			fmt.Printf("构建查询请求失败: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		_, err = conn.Write(packet)
		if err != nil {
			fmt.Printf("发送查询请求失败: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		fmt.Println("已发送状态查询报文")
		time.Sleep(5 * time.Second) // 每5秒发送一次查询
	}
}

// 接收状态发布报文协程
func StartReceiver() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("192.168.0.2"),
		Port: 48000,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 65535)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("接收报文失败: %v\n", err)
			continue
		}

		go handleReceivedPublishPacket(buffer[:n], addr)
	}
}

// 处理接收到的状态发布报文
func handleReceivedPublishPacket(data []byte, addr *net.UDPAddr) {
	header, tlvs, err := canp.ParsePublishPacket(data)
	if err != nil {
		fmt.Printf("解析报文失败: %v\n", err)
		return
	}

	fmt.Printf("=== 收到状态发布报文 ===\n")
	fmt.Printf("源IP: %s\n", addr.IP)
	fmt.Printf("源端口: %d\n", addr.Port)
	fmt.Printf("节点ID: %d\n", header.NodeID)
	fmt.Printf("节点地址: %s\n", header.NodeAddr)
	fmt.Printf("TLV数据:\n")

	for tlvType, tlvValue := range tlvs {
		fmt.Printf("  类型: 0x%02x 长度: %d 值: %v\n", tlvType, tlvValue.Length, tlvValue.Value)
	}
}
