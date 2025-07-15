package main

import (
	"fmt"
	"testing"
	"time"
)

// 测试用例
func TestTokenMonitor(t *testing.T) {

	// 测试参数（来自示例curl命令）
	testURL := "https://api-base.dc-cn.cn.ecouser.net/api/users/newauth/oauth2/token"
	testClientID := "6870ae37a02a3b54c4769c85"
	testClientSecret := "697030a8e56f4af988281b2277cac214"
	testUsername := "kafka_nlp"
	testPassword := "kafka_nlp"

	tm, err := NewTokenManager(testURL, testClientID, testClientSecret, testUsername, testPassword)
	if err != nil {
		t.Fatalf("NewTokenManager error: %v", err)
	}

	// 打印协程：每5秒打印一次（测试用缩短间隔）
	printed := 0
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 安全获取令牌信息并打印
				token, err := tm.GetTokenInfo()
				if err != nil {
					t.Errorf("GetTokenInfo error: %v", err)
					return
				}
				fmt.Printf("Now[%+v]Test Print [%d]: %+v\n", time.Now(), printed+1, token)
				printed++
				if printed >= 3 {
					close(done)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// 等待测试完成
	<-done
	fmt.Println("Test completed")
}
