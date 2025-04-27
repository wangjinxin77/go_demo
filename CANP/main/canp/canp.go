// canp/canp.go
package canp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"net"

	"canp.server/common"
	"canp.server/comp_attr"
)

// CANP报文头部（修正后的24字节结构）
type CANPHeader struct {
	Type     uint8  // CANP 报文类型
	Reserve  uint8  // 保留字段
	NodeID   uint16 // 算力节点唯一标识
	NodeAddr net.IP // 算力节点IPv6地址
	TLVsLen  uint16 // CANP Message Value字段, 即TLV总长度, 单位字节
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
	Header       *CANPHeader
	TLVTypes     []uint8
	SrcIP        net.IP
	SrcPort      uint16
	PeriodicSend bool
	Reserve2     uint8
}

// 创建CANP查询报文（Type=0x02）
func CreateQueryRequest(nodeID uint16, nodeAddr net.IP, tlvTypes []uint8) ([]byte, error) {
	header := createBaseHeader(common.CANPTypeQuery, nodeID, nodeAddr)
	header.TLVsLen = uint16(len(tlvTypes))
	headerBytes := make([]byte, 22)
	headerBytes[0] = header.Type
	headerBytes[1] = header.Reserve
	binary.BigEndian.PutUint16(headerBytes[2:4], header.NodeID)
	copy(headerBytes[4:20], header.NodeAddr.To16())
	binary.BigEndian.PutUint16(headerBytes[20:22], header.TLVsLen)

	return append(headerBytes, tlvTypes...), nil
}

func createBaseHeader(canpType uint8, nodeID uint16, nodeAddr net.IP) *CANPHeader {
	header := &CANPHeader{
		Type:     canpType,
		Reserve:  0x00,
		NodeID:   nodeID,
		NodeAddr: nodeAddr,
	}
	return header
}

// 创建CANP发布报文（Type=0x01）
func CreatePublishPacket(node *comp_attr.NodeInfo) ([]byte, error) {
	header := createBaseHeader(common.CANPTypePublish, node.NodeID, node.NodeAddr)
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
			length := make([]byte, tlvLenSize)
			// 注意，由于代码里tlv.Length是uint16类型，但发的包定义为uint32, 所以这里需要转换为uint32
			binary.BigEndian.PutUint32(length, uint32(tlv.Length))
			tlvBytes = append(tlvBytes, length...)
		} else {
			length := make([]byte, tlvLenSize)
			binary.BigEndian.PutUint16(length, tlv.Length)
			tlvBytes = append(tlvBytes, length...)
		}

		tlvBytes = append(tlvBytes, tlv.Value...)
		totalTLVLength += 2 + tlvLenSize + tlv.Length
		// fmt.Printf("打包发布报文TLV类型: 0x%02x 长度: %d 值: %v, 此时tlv byte长度%d, 值 %v\n", tlv.Type, tlv.Length, tlv.Value, totalTLVLength, tlvBytes)
	}
	header.TLVsLen = totalTLVLength

	headerBytes := make([]byte, common.CANPHederLen)
	headerBytes[0] = header.Type
	// fmt.Printf("打包发布报文NodeID: 数值: 0x%02x byte值: %v\n", header.Type, headerBytes[0])
	headerBytes[1] = header.Reserve
	binary.BigEndian.PutUint16(headerBytes[2:4], header.NodeID)
	// fmt.Printf("打包发布报文NodeID: 数值: %d byte值: %v\n", header.NodeID, headerBytes[2:4])
	copy(headerBytes[4:20], header.NodeAddr.To16())
	// fmt.Printf("打包发布报文NodeAddr: 数值: %v byte值: %v\n", header.NodeAddr, headerBytes[4:20])
	binary.BigEndian.PutUint16(headerBytes[20:22], header.TLVsLen)
	// fmt.Printf("打包发布报文TLV字节数: 数值: %d byte值: %v\n", header.TLVsLen, headerBytes[20:22])

	return append(headerBytes, tlvBytes...), nil
}

func ParseBase(data []byte) (*CANPHeader, []byte, error) {
	if len(data) <= common.CANPHederLen {
		return nil, nil, errors.New("invalid CANP header length")
	}

	header := &CANPHeader{
		Type:    data[0],
		Reserve: data[1],
		NodeID:  binary.BigEndian.Uint16(data[2:4]),
	}
	header.NodeAddr = net.IP(data[4:20])
	header.TLVsLen = binary.BigEndian.Uint16(data[20:common.CANPHederLen])

	return header, data[common.CANPHederLen:], nil
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

		var tlvValueSizeLen uint16 = 2
		var tlvValueSize uint16
		if bytes.Contains(comp_attr.Float32TLVTypes, []uint8{tlvType}) {
			// 处理float32类型的TLV
			tlvValueSizeLen = 4
			tlvValueSize = uint16(binary.BigEndian.Uint32(payload[2 : 2+tlvValueSizeLen]))
		} else {
			tlvValueSize = binary.BigEndian.Uint16(payload[2 : 2+tlvValueSizeLen])
		}
		// fmt.Printf("解析发布报文TLV类型: 0x%02x 长度: %d 值: %v\n", tlvType, tlvValueSize, payload[2+tlvValueSizeLen:2+tlvValueSizeLen+tlvValueSize])

		tlv := comp_attr.TLV{
			Type:    tlvType,
			Reserve: 0x00,
			Length:  tlvValueSize,
		}
		tlv.Value = append(tlv.Value, payload[2+tlvValueSizeLen:2+tlvValueSizeLen+tlvValueSize]...)
		tlvs[tlvType] = tlv

		payload = payload[2+tlvValueSizeLen+tlvValueSize:]
	}

	return header, tlvs, nil
}

// 解析CANP查询报文
func ParseQueryPacket(data []byte) (*CANPHeader, map[uint8]uint8, error) {
	header, payload, err := ParseBase(data)
	if err != nil {
		return nil, nil, err
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
