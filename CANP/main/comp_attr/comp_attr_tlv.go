package comp_attr

import (
	"encoding/binary"
	"fmt"
	"math"
)

// 将uint16转换为字节数组
func uint16ToBytes(v uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	return buf
}

// 将float32转换为字节数组
func float32ToBytes(v float32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, math.Float32bits(v))
	// fmt.Printf("float32转换为字节数组: %f, %v\n", v, buf)
	return buf
}

// 将字节数组转换为float32
func BytesToFloat32(bytes []byte) (float32, error) {
	if len(bytes) != 4 {
		return 0, fmt.Errorf("invalid byte length for float32: expected 4, got %d", len(bytes))
	}
	bits := binary.BigEndian.Uint32(bytes)
	return math.Float32frombits(bits), nil
}

// TLV结构体
type TLV struct {
	Type    uint8
	Reserve uint8
	Length  uint16
	Value   []byte
}

// 创建uint16类型的TLV
func NewUint16TLV(tlvType uint8, value uint16) TLV {
	valueBytes := uint16ToBytes(value)
	return TLV{
		Type:    tlvType,
		Reserve: 0x00,
		Length:  uint16(len(valueBytes)),
		Value:   valueBytes,
	}
}

// 创建float32类型的TLV
func NewFloat32TLV(tlvType uint8, value float32) TLV {
	valueBytes := float32ToBytes(value)
	return TLV{
		Type:    tlvType,
		Reserve: 0x00,
		Length:  uint16(len(valueBytes)),
		Value:   valueBytes,
	}
}

// CPU资源转换为TLV切片
func (r *CPUResource) ToTLVs() []TLV {
	return []TLV{
		NewUint16TLV(0x01, r.Total),
		NewUint16TLV(0x02, r.Free),
		NewFloat32TLV(0x03, r.UtilRate),
	}
}

// GPU资源转换为TLV切片
func (r *GPUResource) ToTLVs() []TLV {
	return []TLV{
		NewUint16TLV(0x04, r.Total),
		NewUint16TLV(0x05, r.Free),
		NewFloat32TLV(0x06, r.UtilRate),
		NewFloat32TLV(0x07, r.MemUtilRate),
	}
}

// 存储资源转换为TLV切片
func (r *StorageResource) ToTLVs() []TLV {
	return []TLV{
		NewFloat32TLV(0x11, r.TotalStorage),
		NewFloat32TLV(0x12, r.FreeStorage),
		NewUint16TLV(0x13, r.ReadIOPS),
		NewUint16TLV(0x14, r.WriteIOPS),
		NewFloat32TLV(0x15, r.TotalMem),
		NewFloat32TLV(0x16, r.FreeMem),
	}
}

// 网络资源转换为TLV切片
func (r *NetworkResource) ToTLVs() []TLV {
	return []TLV{
		NewFloat32TLV(0x21, r.TotalBandwidth),
		NewFloat32TLV(0x22, r.FreeBandwidth),
	}
}

// NodeInfo转换为TLV切片
func (n *NodeInfo) ToTLVs() []TLV {
	var tlvs []TLV
	tlvs = append(tlvs, n.CPU.ToTLVs()...)
	tlvs = append(tlvs, n.GPU.ToTLVs()...)
	tlvs = append(tlvs, n.Storage.ToTLVs()...)
	tlvs = append(tlvs, n.Network.ToTLVs()...)
	return tlvs
}
