package main

import (
	"fmt"
	"sync"
	"time"
)

// 模拟刷新函数（测试时替换）
var RefreshKafkaTokenFunc = func(apiUrl, clientID, clientSecret, refreshToken string) (string, string, string, int, error) {
	// 模拟成功刷新，返回新的令牌和3600秒有效期
	return "new_id", "new_access", "new_refresh", 3600, nil
}

// TokenManager 管理令牌及监控协程
type TokenManager struct {
	mu        sync.RWMutex  // 读写锁保证并发安全
	tokenInfo *KafkaToken   // 令牌信息
	stopChan  chan struct{} // 停止信号通道
}

// NewTokenManager 创建新的TokenManager, 并GetKafkaToken, 同时启动StartTokenMonitor
func NewTokenManager(apiUrl, clientID, clientSecret, username, password string) (*TokenManager, error) {
	tokenManager := &TokenManager{
		tokenInfo: NewKafkaToken(apiUrl, clientID, clientSecret, username, password),
	}
	if err := tokenManager.tokenInfo.GetKafkaToken(); err != nil {
		return nil, err
	}

	// 启动监控协程
	go tokenManager.StartTokenMonitor()

	return tokenManager, nil
}

// StartTokenMonitor 启动令牌监控协程
func (tm *TokenManager) StartTokenMonitor() {
	tm.stopChan = make(chan struct{})
	defer close(tm.stopChan)

	tm.mu.RLock()
	token := tm.tokenInfo
	tm.mu.RUnlock()

	// 情况1：令牌不存在（AccessToken为空）
	if token.AccessToken == "" {
		return
	}

	for {
		// 主循环执行检查
		tm.checkAndRefresh()
	}
}

// checkAndRefresh 检查令牌状态并触发刷新逻辑
func (tm *TokenManager) checkAndRefresh() {
	tm.mu.RLock()
	token := tm.tokenInfo
	tm.mu.RUnlock()

	deadLine := token.GetTime.Add(time.Second * time.Duration(token.ExpiresIn))
	// 提前12小时预刷新
	preRefresh := deadLine.Add(-12 * time.Hour)
	now := time.Now()
	fmt.Printf("now %+v, preRefresh %+v\n", now, preRefresh)

	// 情况2：需要立即刷新
	if now.After(preRefresh) || now.Equal(preRefresh) {
		tm.refreshToken() // 异步执行刷新
		return
	}

	// 情况3：正常情况，设置延迟刷新
	timer := time.NewTimer(preRefresh.Sub(now))
	defer timer.Stop()

	// 等待定时器触发
	<-timer.C
	tm.refreshToken() // 定时触发刷新
}

// refreshToken 执行令牌刷新
func (tm *TokenManager) refreshToken() {
	tm.mu.Lock()
	// 刷新令牌信息（加写锁）
	fmt.Printf(">>>Now [%+v] start to refresh token", time.Now())
	err := tm.tokenInfo.RefreshKafkaToken()
	tm.mu.Unlock()
	if err == nil {
		return // 成功后退出重试
	}

	// 失败等待5分钟
	backoff := time.Duration(5) * time.Minute
	time.Sleep(backoff)
}

// GetTokenInfo 安全获取当前令牌信息（只读）
func (tm *TokenManager) GetTokenInfo() (KafkaToken, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	select {
	case <-tm.stopChan:
		return *tm.tokenInfo, fmt.Errorf("token %+v 监控程序 StartTokenMonitor 已关闭", *tm.tokenInfo)
	default:
		return *tm.tokenInfo, nil // 返回副本避免外部修改
	}
}

// 主函数（示例运行）
// func main() {
// 	// 示例：初始化并启动监控（包含ApiUrl）
// 	initial := KafkaToken{
// 		AccessToken:  "demo_access",
// 		RefreshToken: "demo_refresh",
// 		ExpiresIn:    5000, // 初始有效期5000秒
// 		ClientID:     "demo_client",
// 		ClientSecret: "demo_secret",
// 		Userid:       "demo_user",
// 		ApiUrl:       "http://kafka.api", // 新增ApiUrl字段
// 	}
// 	tm := NewTokenManager()
// 	go tm.StartTokenMonitor()

// 	// 保持主程序运行
// 	select {}
// }
