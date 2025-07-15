// token_fetcher.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// KafkaToken 存储令牌信息及API地址
type KafkaToken struct {
	ApiUrl       string    `json:"api_url"` // Kafka API地址
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"` // 剩余有效时间（秒）
	GetTime      time.Time `json:"_"`
	Userid       string    `json:"userid"`
}

func NewKafkaToken(apiUrl, clientID, clientSecret, username, password string) *KafkaToken {
	return &KafkaToken{
		ApiUrl:       apiUrl,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
	}
}

// TokenResponse 定义API返回的JSON结构（与真实接口完全一致）
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Userid       string `json:"userid"` // 注意JSON字段是小写userid，结构体标签必须匹配
}

// GetKafkaToken 获取Kafka认证Token的核心函数
// 参数顺序与shell命令完全一致：url, client_id, client_secret, username, password
func (token *KafkaToken) GetKafkaToken() error {
	// 构建表单数据（严格匹配shell命令的参数顺序和编码方式）
	form := url.Values{}
	form.Set("grant_type", "password") // 注意使用Set而非Add（确保单值）
	form.Set("client_id", token.ClientID)
	form.Set("client_secret", token.ClientSecret)
	form.Set("username", token.Username)
	form.Set("password", token.Password)

	// 创建POST请求（显式设置Content-Length避免某些服务器拒绝）
	req, err := http.NewRequest(http.MethodPost, token.ApiUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(form.Encode()))) // 显式设置长度

	// 发送请求（添加超时控制）
	client := &http.Client{
		Timeout: 10 * time.Second, // 关键优化：防止长时间阻塞
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码（真实接口返回200表示成功）
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应内容（便于调试）
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(errBody))
	}

	// 解析响应（严格匹配JSON结构）
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 验证关键字段非空（防御性编程）
	if tokenResp.AccessToken == "" || tokenResp.Userid == "" {
		return fmt.Errorf("响应缺少必要字段（access_token或userid为空）")
	}

	token.Userid = tokenResp.Userid
	token.AccessToken = tokenResp.AccessToken
	token.RefreshToken = tokenResp.RefreshToken
	token.ExpiresIn = tokenResp.ExpiresIn
	token.GetTime = time.Now()

	return nil
}

// TokenRefreshResponse 定义OAuth2令牌响应结构体
type TokenRefreshResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// RefreshKafkaToken 刷新Kafka令牌的函数
// 参数：api地址、client_id、client_secret、refresh_token
// 返回：id_token、access_token、新的refresh_token、过期时间（秒）、错误
func (token *KafkaToken) RefreshKafkaToken() error {
	// 构造表单数据
	formRefresh := url.Values{}
	formRefresh.Add("grant_type", "refresh_token")
	formRefresh.Add("client_id", token.ClientID)
	formRefresh.Add("client_secret", token.ClientSecret)
	formRefresh.Add("refresh_token", token.RefreshToken)

	// 创建POST请求
	req, err := http.NewRequest(http.MethodPost, token.ApiUrl, strings.NewReader(formRefresh.Encode()))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 配置带超时的HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应JSON
	var tokenResp TokenRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	token.AccessToken = tokenResp.AccessToken
	token.RefreshToken = tokenResp.RefreshToken
	token.ExpiresIn = tokenResp.ExpiresIn
	token.GetTime = time.Now()

	return nil
}
