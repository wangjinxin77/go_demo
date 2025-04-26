// tests/edge_test.go
package edge_test

import (
	"testing"

	"canp.server/canp"
	"canp.server/edge"
)

func TestHandleCANPRequest(t *testing.T) {
	// 构造模拟请求（查询CPU和GPU）
	request, _ := canp.CreateQueryRequest()

	// 调用边缘云处理函数
	response, err := edge.HandleCANPRequest(request)
	if err != nil {
		t.Fatalf("HandleCANPRequest failed: %v", err)
	}

	// 验证响应类型
	if response[0] != 0x01 {
		t.Errorf("Invalid response type: %x", response[0])
	}

	// 解析响应TLV
	tlvs, err := canp.ParseResponse(response)
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
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
}
