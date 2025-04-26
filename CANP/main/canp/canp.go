// canp/canp.go
package canp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"time"

	"canp.server/comp_attr"
)

// CANP协议常量
const (
	CANPHederLen    = 22    // CANP报文头部长度
	DefaultPort     = 38000 // 天基边缘云监听端口
	NodeID          = 0x0001
	NodeAddress     = "2001:db8::1"
	ResponseTimeout = 3 * time.Second
)

// CANP报文头部（修正后的24字节结构）
type CANPHeader struct {
	Type     uint8    // CANP 报文类型
	Reserve  uint8    // 保留字段2
	NodeID   uint16   // 算力节点唯一标识
	NodeAddr [16]byte // 算力节点IPv6地址
	TLVsLen  uint16   // CANP Message Value字段, 即TLV总长度, 单位字节
}

// TLV结构体
type TLV struct {
	Type    uint8  // 类型标识
	Reserve uint8  // 预留字段
	Length  uint16 // 值字段长度
	Value   []byte // 实际值
}

// CANP查询请求结构体
type QueryRequest struct {
	Header   CANPHeader
	TLVTypes []uint8
	SrcIP    net.IP
	SrcPort  uint16
	Reserve1 uint8
	Reserve2 uint8
}

// 创建CANP查询报文（Type=0x02）
func CreateQueryRequest(srcIP net.IP, srcPort uint16, nodeID uint16, nodeAddr net.IP, tlvTypes []uint8) ([]byte, error) {
	header := createBaseHeader(0x02, nodeID, nodeAddr)

	// 计算TLV总长度
	var tlvBytes []byte
	for _, t := range tlvTypes {
		tlv := TLV{
			Type:    t,
			Reserve: 0x00,
			Length:  0x00, // 实际长度由SerializeValue计算
		}
		tlvBytes = append(tlvBytes, tlv.Type)
		tlvBytes = append(tlvBytes, tlv.Reserve)
		// 预留长度字段空间，后续填充
		tlvBytes = append(tlvBytes, 0x00, 0x00)
	}

	// 计算TLV总长度并填充
	var totalLength uint16 = 0
	for i := 0; i < len(tlvTypes); i++ {
		tv := tlvTypes[i]
		var value []byte
		switch tv {
		case 0x01: // CPU总核数
			value = Uint16ToBytes(10)
			tlvBytes[5+4*i] = uint8(len(value) >> 8)
			tlvBytes[6+4*i] = uint8(len(value) & 0xFF)
		// 其他TLV类型处理...
		default:
			return nil, fmt.Errorf("unsupported TLV type: %x", tv)
		}
		totalLength += uint16(4 + len(value)) // 每个TLV占4字节头+值长度
	}

	header.TLVsLen = totalLength
	headerBytes := make([]byte, 22)
	headerBytes[0] = header.Type
	headerBytes[1] = header.Reserve
	binary.BigEndian.PutUint16(headerBytes[2:4], header.NodeID)
	copy(headerBytes[4:20], header.NodeAddr[:])
	binary.BigEndian.PutUint16(headerBytes[20:22], header.TLVsLen)

	return append(headerBytes, tlvBytes...), nil
}

func createBaseHeader(canpType uint8, nodeID uint16, nodeAddr net.IP) *CANPHeader {
	header := &CANPHeader{
		Type:     canpType,
		Reserve:  0x00,
		NodeID:   nodeID,
		NodeAddr: To16ByteArray(nodeAddr),
	}
	return header
}

// 创建CANP发布报文（Type=0x01）
func CreatePublishPacket(node *comp_attr.NodeInfo, srcIP net.IP, srcPort uint16) ([]byte, error) {
	header := createBaseHeader(0x01, node.NodeID, node.NodeAddr)
	// 构造TLV数据
	tlvs := node.ToTLVs()

	// 计算总长度
	var totalTLVLength uint16
	var tlvBytes []byte
	for _, tlv := range tlvs {
		tlvBytes = append(tlvBytes, tlv.Type)
		tlvBytes = append(tlvBytes, 0x00)

		var tlvLenSize uint16 = 2
		if bytes.Contains(comp_attr.Float32TLVTypes, []uint8{tlv.Type}) {
			// 处理float32类型的TLV
			tlvLenSize = 4
			// 注意，由于代码里tlv.Length是uint16类型，但发的包定义为uint32, 所以这里需要转换为uint32
			binary.BigEndian.PutUint32(tlvBytes[len(tlvBytes):len(tlvBytes)+4], uint32(tlv.Length))
		} else {
			binary.BigEndian.PutUint16(tlvBytes[len(tlvBytes):len(tlvBytes)+2], tlv.Length)
		}

		tlvBytes = append(tlvBytes, tlv.Value...)
		totalTLVLength += 2 + tlvLenSize + tlv.Length
	}
	header.TLVsLen = totalTLVLength

	headerBytes := make([]byte, 24)
	headerBytes[0] = header.Type
	headerBytes[1] = header.Reserve
	binary.BigEndian.PutUint16(headerBytes[2:4], header.NodeID)
	copy(headerBytes[4:20], header.NodeAddr[:])
	binary.BigEndian.PutUint16(headerBytes[20:22], header.TLVsLen)

	return append(headerBytes, tlvBytes...), nil
}

func ParseBase(data []byte) (*CANPHeader, []byte, error) {
	if len(data) <= CANPHederLen {
		return nil, nil, errors.New("invalid CANP header length")
	}

	header := &CANPHeader{
		Type:    data[0],
		Reserve: data[1],
		NodeID:  binary.BigEndian.Uint16(data[2:4]),
	}
	copy(header.NodeAddr[:], data[4:20])
	header.TLVsLen = binary.BigEndian.Uint16(data[20:CANPHederLen])

	return header, data[CANPHederLen:], nil
}

// 解析CANP发布报文
func ParsePublishPacket(data []byte) (*CANPHeader, map[uint8]comp_attr.TLV, error) {
	header, payload, err := ParseBase(data)
	if err != nil {
		return nil, nil, err
	}
	if len(payload) < 4 {
		return nil, nil, errors.New("invalid payload length")
	}
	// 解析TLV数据
	tlvs := make(map[uint8]comp_attr.TLV)
	for len(payload) > 0 {
		if len(payload) < 4 {
			return nil, nil, errors.New("invalid TLV format")
		}
		tlvType := payload[0]
		// tlvReserve := payload[1]
		var tlvLenSize uint16 = 2
		if bytes.Contains(comp_attr.Float32TLVTypes, []uint8{tlvType}) {
			// 处理float32类型的TLV
			tlvLenSize = 4
		}
		tlvLen := binary.BigEndian.Uint16(payload[2 : 2+tlvLenSize])
		tlvs[tlvType] = comp_attr.TLV{
			Type:    tlvType,
			Reserve: 0x00,
			Length:  tlvLen,
			Value:   payload[2+tlvLenSize : 2+tlvLenSize+tlvLen],
		}
		payload = payload[2+tlvLenSize+tlvLen:]
	}

	return header, tlvs, nil
}

// 解析CANP查询报文
func ParseQueryPacket(data []byte) (*CANPHeader, map[uint8]uint8, error) {
	header, payload, err := ParseBase(data)
	if err != nil {
		return nil, nil, err
	}
	if len(payload) < 4 {
		return nil, nil, errors.New("invalid payload length")
	}
	// 解析查询的TLV类型
	tlvs := make(map[uint8]uint8)
	for len(payload) > 0 {
		tlvType := payload[0]
		tlvs[tlvType] = 1 // 表示该TLV类型存在
		payload = payload[1:]
	}

	return header, tlvs, nil
}

// 辅助函数：序列化通用值类型
func Uint16ToBytes(v uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	return buf
}

func Float32ToBytes(v float32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, math.Float32bits(v))
	return buf
}

// IPv6地址转换辅助函数
func To16ByteArray(ip net.IP) [16]byte {
	var addr [16]byte
	copy(addr[:], ip.To16())
	return addr
}
