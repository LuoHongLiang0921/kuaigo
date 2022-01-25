package kafka

import "github.com/Shopify/sarama"

type ClientOption func(c *Client)

// WithConsumerSaramaConfig
// 	@Description 设置kafka 重置kafka 消息
//	@Param cfg sarama.Config
// 	@Return ClientOption 配置项
func WithConsumerSaramaConfig(cfg *sarama.Config) ClientOption {
	return func(c *Client) {
		c.OverwriteConsumerSaramaConfig = cfg
	}
}

// WithProducerSaramaConfig
// 	@Description 设置生产/发布者配置
//	@Param cfg sarama.Config
// 	@Return ClientOption kafka 配置项
func WithProducerSaramaConfig(cfg *sarama.Config) ClientOption {
	return func(c *Client) {
		c.OverwriteProducerSaramaConfig = cfg
	}
}
