package mqmanager

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq/kafka"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq/rabbitmq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq/rocketmq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"time"


	"github.com/mitchellh/mapstructure"
)

// Config 配置
type Config struct {
	// Brokers 消息中间件地址列表
	Brokers []string
	// Mode 队列模式,kafka,rocketmq,rabbitmq
	Mode string
	// RunType 运行类型  发布："publish" 消费："consumer"
	RunType mq.RunType

	Kafka  KafkaConfig  `mapstructure:"kafka"`
	Rabbit RabbitConfig `mapstructure:"rabbit"`
	Rocket RocketConfig `mapstructure:"rocket"`
}

type KafkaConfig struct {
	// KafkaVersion kafka 版本
	Version   string
	Publisher KafkaPublisher `mapstructure:"publisher"`
	Consumer  KafkaConsumer  `mapstructure:"consumer"`
}

type KafkaPublisher struct {
	// 是否异步发布消息
	Async bool
	// 发布消息超时时间
	PublishTimeout time.Duration
	// RequiredAcks 0:NoResponse 1:WaitForLocal -1:WaitForAll
	RequiredAcks int
	Backoff      time.Duration
	// Partitioner random:RandomPartitioner roundRobin:RoundRobinPartitioner hash:HashPartitioner
	Partitioner string
}

type KafkaConsumer struct {
	// Topic 名字
	Topic []string
	// OffsetsInitial offset是否从时间最远开始 1: 最近 2: 最旧
	OffsetsInitial int
	//GroupID 消费组id,如果为空就是按照不是消费者组来走
	GroupID string
	// Assinor  重平衡分配策略 sticky roundRobin range
	Assinor string
}

type RabbitConfig struct {
	// 最大空闲连接数
	MaxIdle int
	// 最大活动连接数
	MaxActive int
	// IdleTimeout 池中空闲时间 秒，每次获取都会清理池中的连接，每次都会从
	IdleTimeout time.Duration
	// MaxConnLifetime 最大存活时间，
	MaxConnLifetime time.Duration
	// 下面试消费者配置

	Consumer RabbitConsumer `mapstructure:"consumer"`
}

type RabbitConsumer struct {
	// Queue 消费
	Queue       string
	ConsumerTag string
	AutoAck     bool
	Exclusive   bool
	NoLocal     bool
	NoWait      bool
	Args        map[string]interface{}

	IsNackRequeue bool
}

type RocketConfig struct {
	GroupName string
	// Topic 名字
	Topic []string

	Retry int
	Async bool
}

// Load
// 	@Description 载入配置，生成多个配置
//	@Param ctx 上下文
// 	@Return map[string]mq.MQ 配合
func Load(ctx context.Context) map[string]mq.MessageQueer {
	var result map[string]mq.MessageQueer
	result = doLoad(ctx)
	conf.OnChange(func(configuration *conf.Configuration) {
		result = doLoad(ctx)
	})
	return result
}

func doLoad(ctx context.Context) map[string]mq.MessageQueer {
	dataV := conf.GetStringMap(constant.ConfigRootKey)
	configs := make(map[string]*Config)
	decoderConfig := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     &configs,
	}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panic("unmarshal key", klog.FieldMod("mq"), klog.FieldErr(err), klog.FieldKey(constant.ConfigRootKey))
		return nil
	}
	err = decoder.Decode(dataV)
	if err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panic("unmarshal key", klog.FieldMod("mq"), klog.FieldErr(err), klog.FieldKey(constant.ConfigRootKey))
		return nil
	}
	result := make(map[string]mq.MessageQueer, len(configs))
	for k, cfg := range configs {
		switch cfg.Mode {
		case constant.ModeKafka:
			result[k] = kafka.Config{
				Mode:    cfg.Mode,
				RunType: cfg.RunType,
				Brokers: cfg.Brokers,
				Version: cfg.Kafka.Version,
				PublishConfig: &kafka.PublishConfig{
					PublishTimeout: cfg.Kafka.Publisher.PublishTimeout,
					Async:          cfg.Kafka.Publisher.Async,
					RequiredAcks:   cfg.Kafka.Publisher.RequiredAcks,
					Backoff:        cfg.Kafka.Publisher.Backoff,
					Partitioner:    cfg.Kafka.Publisher.Partitioner,
				},
				ConsumerConfig: &kafka.ConsumerConfig{
					GroupID:        cfg.Kafka.Consumer.GroupID,
					OffsetsInitial: cfg.Kafka.Consumer.OffsetsInitial,
					Assinor:        cfg.Kafka.Consumer.Assinor,
					Topic:          cfg.Kafka.Consumer.Topic,
				},
			}.Build(ctx)
		case constant.ModeRocketmq:
			result[k] = rocketmq.Config{
				Mode:      cfg.Mode,
				RunType:   cfg.RunType,
				Brokers:   cfg.Brokers,
				Topic:     cfg.Rocket.Topic[0],
				GroupName: cfg.Rocket.GroupName,
				Retry:     cfg.Rocket.Retry,
				Async:     cfg.Rocket.Async,
			}.Build(ctx)
		case constant.ModeRabbitmq:
			result[k] = rabbitmq.Config{
				Address:         cfg.Brokers[0],
				Mode:            cfg.Mode,
				RunType:         cfg.RunType,
				MaxActive:       cfg.Rabbit.MaxActive,
				MaxIdle:         cfg.Rabbit.MaxIdle,
				MaxConnLifetime: cfg.Rabbit.MaxConnLifetime,
				IdleTimeout:     cfg.Rabbit.IdleTimeout,
				ConsumeConfig: rabbitmq.ConsumeConfig{
					Queue:         cfg.Rabbit.Consumer.Queue,
					ConsumerTag:   cfg.Rabbit.Consumer.ConsumerTag,
					AutoAck:       cfg.Rabbit.Consumer.AutoAck,
					Exclusive:     cfg.Rabbit.Consumer.Exclusive,
					NoLocal:       cfg.Rabbit.Consumer.NoLocal,
					NoWait:        cfg.Rabbit.Consumer.NoWait,
					Args:          cfg.Rabbit.Consumer.Args,
					IsNackRequeue: cfg.Rabbit.Consumer.IsNackRequeue,
				},
			}.Build(ctx)
		default:
			klog.TaskLogger.WithContext(ctx).Panicf("message queue mode %s not support", cfg.Mode)
		}
	}
	return result
}
