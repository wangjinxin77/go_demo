package comp_attr

import (
	"net"

	"canp.server/common"
)

// CPU资源结构体
type CPUResource struct {
	Total    uint16  `json:"total"`     // 总量，TLV type=0x01
	Free     uint16  `json:"free"`      // 空闲量，TLV type=0x02
	UtilRate float32 `json:"util_rate"` // 利用率百分比，TLV type=0x03
}

// GPU资源结构体
type GPUResource struct {
	Total       uint16  `json:"total"`         // 总量，TLV type=0x04
	Free        uint16  `json:"free"`          // 空闲量，TLV type=0x05
	UtilRate    float32 `json:"util_rate"`     // 利用率百分比，TLV type=0x06
	MemUtilRate float32 `json:"mem_util_rate"` // GPU内存利用率百分比，TLV type=0x07
}

// 存储资源结构体
type StorageResource struct {
	TotalStorage float32 `json:"total_storage"` // 总存储量(GB)，TLV type=0x11
	FreeStorage  float32 `json:"free_storage"`  // 空闲存储量(GB)，TLV type=0x12
	ReadIOPS     uint16  `json:"read_iops"`     // 读取IOPS，TLV type=0x13
	WriteIOPS    uint16  `json:"write_iops"`    // 写入IOPS，TLV type=0x14
	TotalMem     float32 `json:"total_mem"`     // 总内存(GB)，TLV type=0x15
	FreeMem      float32 `json:"free_mem"`      // 空闲内存(GB)，TLV type=0x16
}

// 网络资源结构体
type NetworkResource struct {
	TotalBandwidth float32 `json:"total_bandwidth"` // 总带宽(Gbps)，TLV type=0x21
	FreeBandwidth  float32 `json:"free_bandwidth"`  // 空闲带宽(Gbps)，TLV type=0x22
}

// float32类型字段的TLV类型常量数组
var Float32TLVTypes = []uint8{
	0x03, // CPUResource.UtilRate
	0x06, // GPUResource.UtilRate
	0x07, // GPUResource.MemUtilRate
	0x11, // StorageResource.TotalStorage
	0x12, // StorageResource.FreeStorage
	0x15, // StorageResource.TotalMem
	0x16, // StorageResource.FreeMem
	0x21, // NetworkResource.TotalBandwidth
	0x22, // NetworkResource.FreeBandwidth
}

// 算力节点信息结构体
type NodeInfo struct {
	NodeID   uint16          `json:"node_id"`   // 节点唯一标识
	NodeAddr net.IP          `json:"node_addr"` // 节点IPv6地址
	CPU      CPUResource     `json:"cpu"`       // CPU资源状态
	GPU      GPUResource     `json:"gpu"`       // GPU资源状态
	Storage  StorageResource `json:"storage"`   // 存储资源状态
	Network  NetworkResource `json:"network"`   // 网络资源状态
}

// 示例节点初始化数据
var DefaultNode = &NodeInfo{
	NodeID:   common.NodeID,
	NodeAddr: net.IPv6loopback,
	CPU: CPUResource{
		Total:    32,
		Free:     24,
		UtilRate: 25.0,
	},
	GPU: GPUResource{
		Total:       8,
		Free:        6,
		UtilRate:    25.0,
		MemUtilRate: 24.0,
	},
	Storage: StorageResource{
		TotalStorage: 1000, // 1TB
		FreeStorage:  800,  // 800GB
		ReadIOPS:     15000,
		WriteIOPS:    12000,
		TotalMem:     128, // 128GB
		FreeMem:      96,  // 96GB
	},
	Network: NetworkResource{
		TotalBandwidth: 1000, // 1Gbps
		FreeBandwidth:  500,  // 500Mbps
	},
}
