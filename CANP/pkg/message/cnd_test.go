package cdn

import (
    "fmt"
    "net"
    "reflect"
    "testing"
)

// 测试EncodeTLV函数
func TestEncodeTLV(t *testing.T) {
    tests := []struct {
        name    string
        tlv     TLV
        want    []byte
        wantErr bool
    }{
        {
            name: "CPU总核数",
            tlv: TLV{
                Type:  CPUCoreTotal,
                Value: uint32(8),
            },
            want: []byte{0x01, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x08},
        },
        {
            name: "GPU利用率",
            tlv: TLV{
                Type:  GPUUtilization,
                Value: float32(75.5),
            },
            want: []byte{0x06, 0x00, 0x00, 0x04, 0x43, 0x16, 0x00, 0x00},
        },
        {
            name: "存储读IOPS",
            tlv: TLV{
                Type:  StorageReadIOPS,
                Value: uint16(1500),
            },
            want: []byte{0x13, 0x00, 0x00, 0x02, 0x05, 0xDC},
        },
        {
            name: "无效类型",
            tlv: TLV{
                Type:  0xFF,
                Value: uint32(100),
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := EncodeTLV(tt.tlv)
            if (err != nil) != tt.wantErr {
                t.Errorf("EncodeTLV() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("EncodeTLV() got = %v, want %v", got, tt.want)
            }
        })
    }
}

// 测试AssembleCANP函数
func TestAssembleCANP(t *testing.T) {
    header := CANPHeader{
        Type:        0x01,
        Reserved:    0x00,
        NodeID:      0x0001,
        NodeAddress: net.ParseIP("2001:db8::1"),
    }

    tests := []struct {
        name    string
        header  CANPHeader
        tlvs    []TLV
        want    []byte
        wantErr bool
    }{
        {
            name: "基本CANP报文",
            header: header,
            tlvs: []TLV{
                {Type: CPUCoreTotal, Value: uint32(8)},
                {Type: GPUCountTotal, Value: uint16(4)},
            },
            want: []byte{
                0x01, 0x00, // 报文类型
                0x00, 0x00, // 预留
                0x00, 0x18, // 总长度=20+8=28
                0x00, 0x01, // 节点ID
                0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // IPv6地址
                0x01, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x08, // CPU总核数TLV
                0x04, 0x00, 0x00, 0x02, 0x00, 0x04, // GPU总数TLV
            },
        },
        {
            name: "错误TLV类型",
            header: header,
            tlvs: []TLV{
                {Type: 0xFF, Value: uint32(100)},
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := AssembleCANP(tt.header, tt.tlvs...)
            if (err != nil) != tt.wantErr {
                t.Errorf("AssembleCANP() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("AssembleCANP() got = %v, want %v", got, tt.want)
            }
        })
    }
}

// 测试调用发包函数
func TestMain(t *testing.T)  {
    // 示例用法
    header := CANPHeader{
        Type:        0x01,
        Reserved:    0x00,
        NodeID:      0x0001,
        NodeAddress: net.ParseIP("2001:db8::1"),
    }

    tlvs := []TLV{
        {Type: CPUCoreTotal, Value: uint32(8)},
        {Type: GPUUtilization, Value: float32(75.5)},
    }

    canpData, err := AssembleCANP(header, tlvs...)
    if err != nil {
        fmt.Println("组装失败:", err)
        return
    }

    fmt.Printf("生成的CANP报文: %x\n", canpData)
}