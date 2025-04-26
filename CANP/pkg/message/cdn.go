package cdn

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "math"
    "net"
)

// TLV类型定义（Type-Length-Value）
const (
    CPUCoreTotal      = 0x01 // CPU总核数
    CPURemaining      = 0x02 // CPU剩余核数
    CPUUtilization    = 0x03 // CPU利用率
    GPUCountTotal     = 0x04 // GPU总数
    GPUCountRemaining = 0x05 // GPU剩余数量
    GPUUtilization    = 0x06 // GPU利用率
    GPUMemoryUtil     = 0x07 // GPU内存利用率
    StorageTotal      = 0x11 // 存储总量
    StorageFree       = 0x12 // 存储剩余量
    StorageReadIOPS   = 0x13 // 存储读IOPS
    StorageWriteIOPS  = 0x14 // 存储写IOPS
    MemoryTotal       = 0x15 // 内存总量
    MemoryAvailable   = 0x16 // 可用内存
    NetBandwidthTotal = 0x21 // 网络总带宽
    NetBandwidthFree  = 0x22 // 网络剩余带宽
)

// TLV结构体定义
type TLV struct {
    Type   uint8       // TLV类型
    Value  interface{} // TLV值，支持多种类型
}

// CANP报文头部结构体定义
type CANPHeader struct {
    Type          uint8  // 报文类型
    Reserved      uint8  // 预留字段
    Length        uint16 // 报文总长度（包括头部和TLV数据）
    NodeID        uint16 // 节点ID
    NodeAddress   net.IP // 节点地址（IPv6）
}

// TLV编码器：将TLV结构体编码为字节数组
func EncodeTLV(t TLV) ([]byte, error) {
    var data []byte

    // 编码类型和预留字段
    header := make([]byte, 2)
    header[0] = t.Type
    header[1] = 0x00 // 预留字段固定为0

    // 编码值部分，根据值的类型进行处理
    var value []byte
    switch v := t.Value.(type) {
    case uint32:
        value = make([]byte, 4)
        binary.BigEndian.PutUint32(value, v)
    case uint16:
        value = make([]byte, 2)
        binary.BigEndian.PutUint16(value, v)
    case float32:
        value = make([]byte, 4)
        binary.BigEndian.PutUint32(value, math.Float32bits(v))
    case string:
        value = []byte(v)
    default:
        return nil, fmt.Errorf("unsupported type: %T", v)
    }

    // 编码长度字段
    length := make([]byte, 2)
    binary.BigEndian.PutUint16(length, uint16(len(value)))

    // 组合成完整的TLV结构
    data = append(data, header...)
    data = append(data, length...)
    data = append(data, value...)

    return data, nil
}

// CANP报文组装器：将头部和多个TLV编码为完整的CANP报文
func AssembleCANP(header CANPHeader, tlvs ...TLV) ([]byte, error) {
    // 编码所有TLV数据
    var tlvData []byte
    for _, t := range tlvs {
        encoded, err := EncodeTLV(t)
        if err != nil {
            return nil, err
        }
        tlvData = append(tlvData, encoded...)
    }

    // 更新头部长度字段（头部固定20字节 + TLV数据长度）
    header.Length = uint16(20 + len(tlvData))
    header.NodeAddress = header.NodeAddress.To16() // 确保地址为IPv6格式

    // 序列化头部
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, header.Type); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, header.Reserved); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, header.Length); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, header.NodeID); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, header.NodeAddress); err != nil {
        return nil, err
    }

    // 合并头部和TLV数据
    return append(buf.Bytes(), tlvData...), nil
}