package main

import (
	"fmt"
	"testing"
)

func TestGetKafkaToken(t *testing.T) {
	// 测试参数（来自示例curl命令）
	testURL := "https://api-base.dc-cn.cn.ecouser.net/api/users/newauth/oauth2/token"
	testClientID := "6870ae37a02a3b54c4769c85"
	testClientSecret := "697030a8e56f4af988281b2277cac214"
	testUsername := "kafka_nlp"
	testPassword := "kafka_nlp"
	token := NewKafkaToken(testURL, testClientID, testClientSecret, testUsername, testPassword)

	// 预期结果（来自示例curl输出）
	expectedUserID := "34685e18-f4e2-415e-9293-a58e69512797"
	expectedExpiresIn := 604800

	// 错误检查
	if err := token.GetKafkaToken(); err != nil {
		t.Fatalf("测试失败: 发生意外错误 - %v", err)
	}

	// 结果验证
	if token.Userid != expectedUserID {
		t.Errorf("用户ID不匹配\n期望: %q\n实际: %q", expectedUserID, token.Userid)
	}
	if token.ExpiresIn != expectedExpiresIn {
		t.Errorf("过期时间不匹配\n期望: %d\n实际: %d", expectedExpiresIn, token.ExpiresIn)
	}

	fmt.Printf("token info %+v", *token)
}

func TestRefreshKafkaToken(t *testing.T) {
	// 测试参数（来自示例curl命令）
	testURL := "https://api-base.dc-cn.cn.ecouser.net/api/users/newauth/oauth2/token"
	testClientID := "6870ae37a02a3b54c4769c85"
	testClientSecret := "697030a8e56f4af988281b2277cac214"
	testRefreshToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjIjoiNjg3MGFlMzdhMDJhM2I1NGM0NzY5Yzg1IiwidSI6IjM0Njg1ZTE4LWY0ZTItNDE1ZS05MjkzLWE1OGU2OTUxMjc5NyIsInIiOiJXa2N3ekkiLCJ0IjoiciIsImlhdCI6MTc1MjQ4NDY3NiwiZXhwIjoxNzU1MDc2Njc2fQ.Uqb8IHmoekxaOc5fxT2mf4eVDEi81NpAZXYu9m7kf8U4vpgj6J4bMLuMssiNmdG87hErmrLpLcRhjjncS2jb5PY37mJvBrV2kkj_K9IiYZYiPCY6NAy8InFikkfQmokX9fWFcemmp_RK_QfioTLyqZE9MOwP4QbrHy1EiVxC7aI"

	token := &KafkaToken{
		ApiUrl:       testURL,
		ClientID:     testClientID,
		ClientSecret: testClientSecret,
		RefreshToken: testRefreshToken,
	}
	// 预期结果（来自示例curl输出）
	expectedExpiresIn := 604800

	// 检查错误
	if err := token.RefreshKafkaToken(); err != nil {
		t.Fatalf("函数执行失败: %v", err)
	}

	// 验证结果
	// if idToken != expectedIDToken {
	// 	t.Errorf("IDToken不匹配\n实际: %q\n预期: %q", idToken, expectedIDToken)
	// }
	if token.ExpiresIn != expectedExpiresIn {
		t.Errorf("ExpiresIn不匹配\n实际: %d\n预期: %d", token.ExpiresIn, expectedExpiresIn)
	}
	fmt.Printf("token info %+v", *token)
}
