// edge/edge.go
package edge

import (
	"fmt"
	"net"

	"canp.server/canp"
	"canp.server/comp_attr"
)

// 算力节点信息存储
var nodeDB = map[uint16]*comp_attr.NodeInfo{
	1024: comp_attr.DefaultNode,
}

// 资源查询通道
var queryChan = make(chan *canp.QueryRequest, 100)

// 处理接收到的状态发布报文
func handleReceivedPacket(conn *net.UDPConn, data []byte, addr *net.UDPAddr) {
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

	for tlvType, tlv := range tlvs {
		fmt.Printf("  类型: 0x%02x 长度: %d 值: %v\n", tlvType, tlv.Length, tlv.Value)
	}
}

// 处理接收的CANP查询报文
func StartReceiver() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("192.168.0.1"),
		Port: 38000,
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

		go handleReceivedPacket(conn, buffer[:n], addr)
	}
}

// 其他辅助函数...
