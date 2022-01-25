package kafka

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"time"


	"github.com/Shopify/sarama"
)

const (
	AssinorRange      = "range"
	AssinorRoundRobin = "roundRobin"
	AssinorSticky     = "sticky"

	NoResponse   = 0
	WaitForLocal = 1
	WaitForAll   = -1

	OffsetNewest = 1
	OffsetOldest = 2

	RandomPartitioner     = "random"
	RoundRobinPartitioner = "roundRobin"
	HashPartitioner       = "hash"
)

// Config 配置
type Config struct {
	// Mode 队列模式,kafka,rocketmq,rabbitmq
	Mode string
	// RunType
	RunType mq.RunType
	// Brokers 连接地址列表
	Brokers []string
	// Version kafka 版本
	Version        string
	PublishConfig  *PublishConfig
	ConsumerConfig *ConsumerConfig
}

type PublishConfig struct {
	// 发布消息超时时间
	PublishTimeout time.Duration
	Async          bool
	// RequiredAcks 0:NoResponse 1:WaitForLocal -1:WaitForAll
	RequiredAcks int
	Backoff      time.Duration
	// Partitioner radom:RandomPartitioner RoundRobinPartitioner and HashPartitioner
	Partitioner string
}

type ConsumerConfig struct {
	Topic []string
	//GroupID 消费组id
	GroupID string
	// IsOldest offset是否从时间最远开始 1: 最近 2: 最老
	OffsetsInitial int
	// Assinor  重平衡分配策略 sticky roundrobin range
	Assinor string
}

// RawConfig
// 	@Description  kafka 配置
//	@param ctx 上下文
//	@param key
// 	@return *Config
func RawConfig(ctx context.Context, key string) *Config {
	cfg := getDefaultConfig()
	if err := conf.UnmarshalKey(key, &cfg); err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panicf("kafka config err:%v", err)
	}
	return &cfg
}

func getDefaultConfig() Config {
	return Config{
		Mode: constant.ModeKafka,
	}
}

// Build
// 	@Description 实例化 kafka 实例
// 	@Receiver c Config
//	@Param ctx 上下文
// 	@Return *Client 配置后的kafka 实例
func (c Config) Build(ctx context.Context) *Client {
	c.mustValidConfig()
	c.setDefaults()
	var opts []ClientOption
	if c.RunType.IsConsumerType() {
		opts = c.setConsumerConfig(opts)
	}
	if c.RunType.IsPublishType() {
		opts = c.setPublishConfig(opts)
	}
	client := NewClient(ctx, &c, opts...)
	return client
}

func (c Config) mustValidConfig() {
	if c.Mode == "" {
		klog.KuaigoLogger.Panic("config mode must not empty")
	}
	if c.RunType == "" {
		klog.KuaigoLogger.Panic("config run type must not empty")
	}
	if len(c.Brokers) <= 0 {
		klog.KuaigoLogger.Panic("config brokers must not empty")
	}
	if c.RunType.IsConsumerType() {
		if len(c.ConsumerConfig.Topic) <= 0 {
			klog.KuaigoLogger.Panicf("kafka consumer topic must not empty,but %v ", c.ConsumerConfig.Topic)
		}
	}
}

func (c Config) setPublishConfig(opts []ClientOption) []ClientOption {
	pCfg := sarama.NewConfig()
	pCfg.ClientID = "tabby-kafka"
	pCfg.Version = sarama.V1_0_0_0
	if c.Version != "" {
		kv, err := sarama.ParseKafkaVersion(c.Version)
		if err != nil {
			klog.KuaigoLogger.Panicf("kafka version %v err：%v", c.Version, err)
		}
		pCfg.Version = kv
	}
	pCfg.Producer.Retry.Backoff = c.PublishConfig.Backoff
	pCfg.Producer.Return.Successes = true
	switch c.PublishConfig.Partitioner {
	case HashPartitioner:
		pCfg.Producer.Partitioner = sarama.NewHashPartitioner
	case RandomPartitioner:
		pCfg.Producer.Partitioner = sarama.NewRandomPartitioner
	case RoundRobinPartitioner:
		pCfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	}
	switch c.PublishConfig.RequiredAcks {
	case WaitForAll:
		pCfg.Producer.RequiredAcks = sarama.WaitForAll
	case WaitForLocal:
		pCfg.Producer.RequiredAcks = sarama.WaitForLocal
	case NoResponse:
		pCfg.Producer.RequiredAcks = sarama.NoResponse
	}
	opts = append(opts, WithProducerSaramaConfig(pCfg))
	return opts
}

func (c Config) setConsumerConfig(opts []ClientOption) []ClientOption {
	kCfg := sarama.NewConfig()
	kCfg.ClientID = "tabby-kafka"
	kCfg.Version = sarama.V1_0_0_0
	kCfg.Consumer.Return.Errors = true
	if c.Version != "" {
		kv, err := sarama.ParseKafkaVersion(c.Version)
		if err != nil {
			klog.KuaigoLogger.Panicf("kafka version %v err：%v", c.Version, err)
		}
		kCfg.Version = kv
	}
	switch c.ConsumerConfig.OffsetsInitial {
	case OffsetNewest:
		kCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	case OffsetOldest:
		kCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	default:
		kCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	switch c.ConsumerConfig.Assinor {
	case AssinorRange:
		kCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	case AssinorRoundRobin:
		kCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case AssinorSticky:
		kCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	default:
		kCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	}
	opts = append(opts, WithConsumerSaramaConfig(kCfg))
	return opts
}

func (c *Config) setDefaults() {
	if c.RunType.IsConsumerType() {
		if c.ConsumerConfig == nil {
			klog.KuaigoLogger.Panicf("mq run type %v,but consumer is nil", c.RunType)
		}
		if c.ConsumerConfig.Assinor == "" {
			c.ConsumerConfig.Assinor = AssinorRange
		}
	}

	if c.RunType.IsPublishType() {
		if c.PublishConfig == nil {
			klog.KuaigoLogger.Panicf("mq run type %v,but publish field is nil", c.RunType)
		}
		if c.PublishConfig.PublishTimeout == 0 {
			c.PublishConfig.PublishTimeout = ktime.Duration("5s")
		}
		if c.PublishConfig.Partitioner == "" {
			c.PublishConfig.Partitioner = HashPartitioner
		}

		if c.PublishConfig.Backoff == 0 {
			c.PublishConfig.Backoff = ktime.Duration("2s")
		}
	}

}

func defaultSaramaConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = "tabby-kafka"
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	return config
}

func defaultSaramaSyncPublisherConfig() *sarama.Config {
	pCfg := sarama.NewConfig()
	pCfg.ClientID = "tabby-kafka"
	pCfg.Version = sarama.V1_0_0_0
	pCfg.Metadata.Retry.Backoff = time.Second * 2

	pCfg.Producer.RequiredAcks = sarama.WaitForAll
	pCfg.Producer.Return.Successes = true

	pCfg.Producer.Partitioner = sarama.NewRandomPartitioner
	return pCfg
}

func defaultSaramaASyncPublisherConfig() *sarama.Config {

	pCfg := sarama.NewConfig()
	pCfg.ClientID = "tabby-kafka"
	pCfg.Version = sarama.V1_0_0_0
	pCfg.Metadata.Retry.Backoff = time.Second * 2

	// 等待服务器所有副本都保存成功后的响应
	pCfg.Producer.RequiredAcks = sarama.WaitForAll
	pCfg.Producer.Return.Successes = true
	// 随机向partition发送消息
	pCfg.Producer.Partitioner = sarama.NewRandomPartitioner

	return pCfg
}
