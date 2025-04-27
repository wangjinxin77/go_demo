package edge

import (
	"context"
	"fmt"
	"net"
	"time"

	"canp.server/canp"
	"canp.server/common"
	"canp.server/comp_attr"
)

// 算力节点信息存储
var nodeDB = map[uint16]*comp_attr.NodeInfo{
	common.NodeID: comp_attr.DefaultNode,
}

// 解析接收到的查询报文
func handleReceivedPacket(data []byte, addr *net.UDPAddr) (*canp.QueryRequest, error) {
	header, tlvs, err := canp.ParseQueryPacket(data)
	if err != nil {
		fmt.Printf("解析查询报文失败: %v\n", err)
		return nil, err
	}

	if header.Type == common.CANPTypePublish {
		fmt.Printf("这里不应该收到状态发布报文 - NodeID:%d Addr:%s, 丢弃\n", header.NodeID, header.NodeAddr)
		return nil, nil
	}

	// 提取查询的TLV类型
	var types []uint8
	for typeByte := range tlvs {
		types = append(types, typeByte)
	}

	// 构造查询请求对象
	return &canp.QueryRequest{
		Header:   header,
		SrcIP:    addr.IP,
		SrcPort:  uint16(addr.Port),
		TLVTypes: types,
	}, nil
}

func StartRcvSendPublishPacket(listenIp string, listenPort int, queryChan chan *canp.QueryRequest, ctx context.Context) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(listenIp),
		Port: listenPort,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go startReceiver(conn, queryChan, ctx)
	go startResponder(conn, queryChan, ctx)
	fmt.Println("edge服务已启动")
	// 等待停止信号

	<-ctx.Done()
	fmt.Println("停止edge服务")
}

// 处理接收的CANP查询报文
func startReceiver(conn *net.UDPConn, queryChan chan *canp.QueryRequest, ctx context.Context) {
	buffer := make([]byte, 65535)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("edge接收到结束通知, 停止接收")
			return
		default:

			conn.SetReadDeadline(time.Now().Add(common.ReadTimeout)) // 1秒超时
			n, addr, err := conn.ReadFromUDP(buffer)
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 处理超时逻辑, 超时不做任何操作
				continue
			} else if err != nil {
				fmt.Printf("edge接收报文失败: %v\n", err)
				continue
			}

			req, err := handleReceivedPacket(buffer[:n], addr)
			if err == nil && req != nil {
				fmt.Printf("=== edge收到查询请求 ===\n")
				fmt.Printf("源IP: %s\n", addr.IP)
				fmt.Printf("源端口: %d\n", addr.Port)
				fmt.Printf("节点ID: %d\n", req.Header.NodeID)
				fmt.Printf("节点地址: %s\n", req.Header.NodeAddr)
				fmt.Printf("查询的TLV类型: %v\n", req.TLVTypes)
				fmt.Printf("查询的TLV长度: %d\n", req.Header.TLVsLen)
				// 放入处理通道
				queryChan <- req
			}

		}
	}
}

// 根据地址查找节点
func nodeDBByAddr(addr net.IP) *comp_attr.NodeInfo { // 修改返回类型为*comp_attr.NodeInfo
	for _, node := range nodeDB {
		if node.NodeAddr.Equal(addr) {
			return node
		}
	}
	return nil
}

// 发送CANP发布报文
func sendResponse(conn *net.UDPConn, packet []byte, dstIP net.IP, dstPort uint16) {

	// 解析目标地址
	target, err := net.ResolveUDPAddr("udp", dstIP.String()+fmt.Sprintf(":%d", dstPort))
	if err != nil {
		fmt.Printf("edge地址解析失败: %v\n", err)
	}

	// 通过同一个连接发送
	if _, err := conn.WriteToUDP(packet, target); err != nil {
		fmt.Printf("edge发送失败: %v\n", err)
	} else {
		fmt.Printf("edge[发送] 到 %s:%d 成功\n", target.IP, target.Port)
	}
}

// 处理查询请求并发送响应
func startResponder(conn *net.UDPConn, queryChan chan *canp.QueryRequest, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("edge接收到结束通知, edge停止响应")
			return
		case req := <-queryChan:
			node, exists := nodeDB[req.Header.NodeID]
			if !exists && req.Header.NodeAddr != nil {
				node = nodeDBByAddr(req.Header.NodeAddr)
			}

			if node == nil {
				fmt.Printf("未找到节点: ID=%d Addr=%s\n", req.Header.NodeID, req.Header.NodeAddr)
				continue
			}

			resp, err := canp.CreatePublishPacket(node)
			if err != nil {
				fmt.Printf("构建响应失败: %v\n", err)
				continue
			}

			sendResponse(conn, resp, req.SrcIP, req.SrcPort)
		}
	}
}
