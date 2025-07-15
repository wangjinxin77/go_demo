package entry

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

var config *sarama.Config

// 实现sarama.TokenProvider接口
type OAuthTokenProvider struct {
	tokenManager string
}

func (p *OAuthTokenProvider) Token() (*sarama.AccessToken, error) {
	return &sarama.AccessToken{
		Token: p.tokenManager,
	}, nil
}

func init() {
	config = sarama.NewConfig()
	config.Version = sarama.V0_10_2_1
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Net.DialTimeout = 60 * time.Second
	config.Net.ReadTimeout = 60 * time.Second
	config.Net.WriteTimeout = 60 * time.Second
	config.ClientID = "biz-kafka-consumer"

	config.Version = sarama.V2_6_0_0
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
	config.Net.SASL.TokenProvider = &OAuthTokenProvider{tokenManager: "access_token"}
}

func startConsumerGroup(kafkaUriString string, kafkaType string, handler sarama.ConsumerGroupHandler) {
	groupName, kafkaBroker := "group_id", "topics"
	_, err := sarama.NewConsumerGroup([]string{kafkaBroker}, groupName, config)
	if err != nil {
		log.Fatalf("创建消费者组失败: %v", err)
	}
}
