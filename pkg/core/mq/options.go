package mq

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"strings"
)

// PublishOptions 配置项
type PublishOptions struct {
	RabbitPublishOptions
}

type RabbitPublishOptions struct {
	Exchange        string
	ExchangeType    string
	Reliable        bool
	IsCheckExchange bool
	Durable         bool
}

// PublishOption 发布配置函数
type PublishOption func(*PublishOptions)

// WithRabbitPublishExchange
// 	@Description 设置交换机
//	@Param exchange 交换机名字
// 	@Return PublishOption 设置交换机后的名字
func WithRabbitPublishExchange(exchange string) PublishOption {
	return func(options *PublishOptions) {
		options.Exchange = exchange
	}
}

// WithRabbitPublishExchangeType
// 	@Description
//	@Param exchangeType
// 	@Return PublishOption
func WithRabbitPublishExchangeType(exchangeType string) PublishOption {
	return func(options *PublishOptions) {
		options.ExchangeType = exchangeType
	}
}

// WithRabbitPublishIsCheckExchange
// 	@Description 设置是否检查交换机权限
//	@Param isCheckExchange
// 	@Return PublishOption
func WithRabbitPublishIsCheckExchange(isCheckExchange bool) PublishOption {
	return func(options *PublishOptions) {
		options.IsCheckExchange = isCheckExchange
	}
}

// WithRabbitPublishDurable
// 	@Description 设置是否可持久化
//	@Param durable
// 	@Return PublishOption
func WithRabbitPublishDurable(durable bool) PublishOption {
	return func(options *PublishOptions) {
		options.Durable = durable
	}
}

// WithRabbitPublishReliable
// 	@Description 设置发布是否可以可用性，有返回
//	@Param reliable
// 	@Return PublishOption
func WithRabbitPublishReliable(reliable bool) PublishOption {
	return func(options *PublishOptions) {
		options.Reliable = reliable
	}
}

// ConsumerOptions 消费配置项
type ConsumerOptions struct {
	RabbitConsumerOptions
	KafkaConsumerOptions
}

type RabbitConsumerOptions struct {
	// Rabbit 使用
	Queue string
	// ConsumerTag ...
	ConsumerTag string
	AutoAck     bool
	Exclusive   bool
	NoLocal     bool
	NoWait      bool
	//
	Args map[string]interface{}
}

type KafkaConsumerOptions struct {
}

// ConsumerOption 消费配置项函数
type ConsumerOption func(options *ConsumerOptions)

// WithConsumerQueue
// 	@Description 设置队列名
//	@Param queue
// 	@Return ConsumerOption
func WithConsumerQueue(queue string) ConsumerOption {
	return func(options *ConsumerOptions) {
		options.Queue = queue
	}
}

// MergeConsumerOption
// 	@Description
//	@Param opts
// 	@Return *ConsumerOptions
func MergeConsumerOption(opts ...ConsumerOption) *ConsumerOptions {
	copts := &ConsumerOptions{}
	for _, f := range opts {
		f(copts)
	}
	return copts
}

// MergePublishOption
// 	@Description
//	@Param opts
// 	@Return *PublishOptions
func MergePublishOption(opts ...PublishOption) *PublishOptions {
	copts := &PublishOptions{}
	for _, f := range opts {
		f(copts)
	}
	return copts
}

type RunType string

// IsPublishType
// 	@Description 是否是发布类型
// 	@Receiver r
// 	@Return bool
func (r RunType) IsPublishType() bool {
	if r == "" {
		return false
	}
	if r == constant.RunTypePublish {
		return true
	}
	if strings.Contains(string(r), constant.RunTypePublish) {
		return true
	}
	return false
}

// IsConsumerType
// 	@Description 是否是消费类型
// 	@Receiver r
// 	@Return bool
func (r RunType) IsConsumerType() bool {
	if r == "" {
		return false
	}
	if r == constant.RunTypeConsumer {
		return true
	}
	if strings.Contains(string(r), constant.RunTypeConsumer) {
		return true
	}
	return false
}
