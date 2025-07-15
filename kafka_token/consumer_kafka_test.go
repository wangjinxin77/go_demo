package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"

	"github.com/IBM/sarama"
)

// 配置参数
const (
	kafkaBroker   = "ksync-ngiot-cn.dc.cn.ecouser.net:20000"
	groupName     = "oauth-secure-group"
	tokenEndpoint = "https://auth.example.com/oauth2/token"
)

// 实现sarama.TokenProvider接口
type OAuthTokenProvider struct {
	tokenManager *TokenManager
}

func (p *OAuthTokenProvider) Token() (*sarama.AccessToken, error) {
	token, err := p.tokenManager.GetTokenInfo()
	if err != nil {
		return nil, err
	}
	return &sarama.AccessToken{
		Token:      token.AccessToken,
		Extensions: map[string]string{"client_id": p.tokenManager.tokenInfo.ClientID},
	}, nil
}

// 消费者组处理器
type SecureConsumerGroupHandler struct{}

func (SecureConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (SecureConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h SecureConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Partition: %d | Offset: %d | Key: %s | Value: %s\n",
			msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

		sess.MarkMessage(msg, "")
	}
	return nil
}

func TestKafkaConsumer(t *testing.T) {

	// 测试参数（来自示例curl命令）
	TokenURL := "https://api-base.dc-cn.cn.ecouser.net/api/users/newauth/oauth2/token"
	TokenClientID := "6870ae37a02a3b54c4769c85"
	TokenClientSecret := "697030a8e56f4af988281b2277cac214"
	TokenUsername := "kafka_nlp"
	TokenPassword := "kafka_nlp"

	tm, err := NewTokenManager(TokenURL, TokenClientID, TokenClientSecret, TokenUsername, TokenPassword)
	if err != nil {
		log.Fatalf("NewTokenManager error: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
	config.Net.SASL.TokenProvider = &OAuthTokenProvider{tokenManager: tm}
	config.ClientID = TokenClientID

	//config.ClientID = "k8s-consumer-client"

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup([]string{kafkaBroker}, groupName, config)
	if err != nil {
		log.Fatalf("创建消费者组失败: %v", err)
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Printf("关闭消费者组错误: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	handler := SecureConsumerGroupHandler{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	topics := []string{"IOT2.RCP.0gxgac.onMapSet_V2", "IOT2.RCP.egq7t6.onMapSet_V2"}

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Consume(ctx, topics, handler); err != nil {
					log.Printf("消费错误: %v", err)
				}
				if ctx.Err() != nil {
					return
				}
			}
		}
	}()

	<-signals
	cancel()
	wg.Wait()
}
