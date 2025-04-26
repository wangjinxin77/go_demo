// tests/carrier_test.go
package carrier_test

import (
	"testing"

	"canp.server/carrier"
)

func TestQueryEdgeNode(t *testing.T) {
	// 使用模拟的边缘云地址
	nodeAddr := "2001:db8::1"

	// 发起查询请求
	tlvs, err := carrier.QueryEdgeNode(nodeAddr)
	if err != nil {
		t.Fatalf("QueryEdgeNode failed: %v", err)
	}

	// 验证必须包含的TLV类型
	expectedTypes := map[uint8]bool{
		0x01: true, // CPU总核数
		0x04: true, // GPU总个数
		0x05: true, // GPU剩余个数
	}

	for tlvType := range tlvs {
		if !expectedTypes[tlvType] {
			t.Errorf("Unexpected TLV type: %x", tlvType)
		}
		delete(expectedTypes, tlvType)
	}

	if len(expectedTypes) > 0 {
		t.Errorf("Missing expected TLV types: %v", expectedTypes)
	}

	// 验证TLV值类型
	value := tlvs[0x01]  // 正确：直接使用 []byte
	cpuCores := value[0] // Get first byte as uint8
	if cpuCores == 0 {
		t.Error("CPU总核数不应为0")
	}
}
