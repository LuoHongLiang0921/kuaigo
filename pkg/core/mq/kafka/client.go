package kafka

import (
	"context"
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kgo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Client struct {
	cfg *Config

	ctx context.Context
	// 消费使用
	//OverwriteConsumerSaramaConfig sarama 消费配置
	OverwriteConsumerSaramaConfig *sarama.Config
	consumerGroup                 sarama.ConsumerGroup
	// wg，等待消费消息
	wg sync.WaitGroup

	//OverwriteProducerSaramaConfig sarama 生产配置
	OverwriteProducerSaramaConfig *sarama.Config
	// asyncProducer 异步生产者
	asyncProducer sarama.AsyncProducer
	// syncProducer 同步生产者
	syncProducer sarama.SyncProducer

	// 下面为控制信号
	//stopConsumerChan chan struct{}
	closing chan struct{}
	// 业务处理函数
	messageHandlerFunc mq.HandlerFunc
	logger             *klog.Logger
}

// NewClient
// 	@Description  new client
//	@Param cfg 配置
// 	@Return *Client 初始化资源和
func NewClient(ctx context.Context, cfg *Config, opts ...ClientOption) *Client {
	c := &Client{
		cfg: cfg,
	}
	for _, f := range opts {
		f(c)
	}
	c.closing = make(chan struct{})
	c.ctx = ctx
	if cfg.RunType.IsConsumerType() {
		c.newConsumerGroup(ctx)
	}
	if cfg.RunType.IsPublishType() {
		c.newProducer(ctx)
	}
	c.logger = klog.KuaigoLogger
	return c
}

// Publish
// 	@Description 发布消息,如果 配置中async =true
// 	@Receiver c Client
//	@Param ctx 上下文
//	@Param topic 主题
//	@Param msg 消息内容
// 	@Return error 错误
func (c *Client) Publish(ctx context.Context, topic string, msg *mq.Message, opts ...mq.PublishOption) (*mq.RespMessage, error) {
	if c.cfg.PublishConfig.Async {
		return c.doAsyncPublish(ctx, topic, msg)
	}
	return c.doSyncPublish(ctx, topic, msg)
}

func (c *Client) doSyncPublish(ctx context.Context, topic string, msg *mq.Message) (*mq.RespMessage, error) {
	pMsg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(msg.Body),
		Timestamp: time.Now(),
	}
	p, offset, err := c.syncProducer.SendMessage(pMsg)
	if err != nil {
		return nil, err
	}
	return &mq.RespMessage{
		Topic:     topic,
		Offset:    offset,
		Partition: p,
	}, nil
}

func (c *Client) doAsyncPublish(ctx context.Context, topic string, msg *mq.Message) (*mq.RespMessage, error) {
	pMsg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(msg.Body),
		Timestamp: time.Now(),
	}
	c.asyncProducer.Input() <- pMsg
	ctx, cancel := context.WithTimeout(ctx, c.cfg.PublishConfig.PublishTimeout)
	defer cancel()

	select {
	case info := <-c.asyncProducer.Successes():
		return &mq.RespMessage{
			Topic:     info.Topic,
			Offset:    info.Offset,
			Partition: info.Partition,
			Timestamp: info.Timestamp,
		}, nil
	case fail := <-c.asyncProducer.Errors():
		if nil != fail {
			return nil, fail.Err
		}
	case <-ctx.Done():
		return nil, errors.New("publish message timeout")
	}
	return nil, nil
}

// RegisterHandler
// 	@Description  业务端需要实现的函数
// 	@Receiver c Client
//	@Param ctx 上下文
//	@Param h 处理函数
// 	@Return *Client 设置处理函数后的消费实例
func (c *Client) RegisterHandler(ctx context.Context, h mq.HandlerFunc) error {
	c.messageHandlerFunc = h
	return nil
}

// Consume
// 	@Description  消费消息，消息需要 mq.Message  中的 Ack() 或 NAck() 方法，否则会阻塞
// 	@Receiver c Client
// 	@Return error 错误
func (c *Client) Consume(ctx context.Context, opts ...mq.ConsumerOption) error {
	c.logger.WithContext(ctx).Infof("start kafka task,topic :%v ", c.cfg.ConsumerConfig.Topic[0])

	msgs, err := c.Subscribe(ctx, c.cfg.ConsumerConfig.Topic[0])
	if err != nil {
		return err
	}
	for msg := range msgs {
		_ = c.messageHandlerFunc(ctx, msg)
	}
	return nil
}

// Stop
// 	@Description 停止
// 	@Receiver c Client
// 	@Return error 错误
func (c *Client) Stop() error {
	close(c.closing)
	c.wg.Wait()
	return nil
}

func (c *Client) GracefulStop() error {
	close(c.closing)
	c.wg.Wait()
	return nil
}

func (c *Client) newConsumerGroup(ctx context.Context) *Client {
	if c.OverwriteConsumerSaramaConfig == nil {
		c.OverwriteConsumerSaramaConfig = defaultSaramaConsumerConfig()
	}
	if c.cfg.Version != "" {
		kv, err := sarama.ParseKafkaVersion(c.cfg.Version)
		if err != nil {
			c.logger.Panicf("kafka version %v err：%v", c.cfg.Version, err)
		}
		c.OverwriteConsumerSaramaConfig.Version = kv
	}
	csg, err := sarama.NewConsumerGroup(c.cfg.Brokers, c.cfg.ConsumerConfig.GroupID, c.OverwriteConsumerSaramaConfig)
	if err != nil {
		c.logger.WithContext(ctx).Panicf("new consumer group err:%v", err)
	}
	kgo.Go(func() {
		for err := range csg.Errors() {
			if err == nil {
				continue
			}
			c.logger.Errorf("sarama internal error %v", err)
		}
	})
	c.consumerGroup = csg
	return c
}

func (c *Client) newProducer(ctx context.Context) *Client {
	if c.OverwriteProducerSaramaConfig == nil {
		if c.cfg.PublishConfig.Async {
			c.OverwriteProducerSaramaConfig = defaultSaramaASyncPublisherConfig()
		} else {
			c.OverwriteProducerSaramaConfig = defaultSaramaSyncPublisherConfig()
		}
	}
	var err error
	if c.cfg.PublishConfig.Async {
		c.asyncProducer, err = sarama.NewAsyncProducer(c.cfg.Brokers, c.OverwriteProducerSaramaConfig)
	} else {
		c.syncProducer, err = sarama.NewSyncProducer(c.cfg.Brokers, c.OverwriteProducerSaramaConfig)
	}

	if err != nil {
		c.logger.WithContext(ctx).Panicf(err.Error())
	}
	return c
}
