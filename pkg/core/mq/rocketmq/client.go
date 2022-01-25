package rocketmq

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kgo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/apache/rocketmq-client-go/v2"
)

type Client struct {
	cfg         *Config
	producerIns rocketmq.Producer
	consumerIns rocketmq.PushConsumer
	Handler     mq.HandlerFunc
	wg          sync.WaitGroup
}

// NewClient
// 	@Description  new client
//	@Param cfg 配置
// 	@Return *Client 初始化资源和
func NewClient(ctx context.Context, cfg *Config, opts ...ClientOption) *Client {
	c := &Client{cfg: cfg}
	for _, f := range opts {
		f(c)
	}
	if cfg.RunType.IsConsumerType() {
		c.newConsumer(ctx)
	}
	if cfg.RunType.IsPublishType() {
		c.newProducer(ctx)
	}
	return c
}

// Consume
// 	@Description 消费
// 	@Receiver c Client
//	@Param ctx 上下文
// 	@Return error 错误
func (c *Client) Consume(ctx context.Context, opts ...mq.ConsumerOption) error {
	err := c.consumerIns.Subscribe(c.cfg.Topic, consumer.MessageSelector{}, c.doConsumer)
	if err != nil {
		return err
	}
	err = c.consumerIns.Start()
	if err != nil {
		_ = c.consumerIns.Unsubscribe(c.cfg.Topic)
		return err
	}
	return nil
}

func (c *Client) doConsumer(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, item := range msgs {
		if c.cfg.Async {
			c.wg.Add(1)
			kgo.SafeGo(func() {
				defer c.wg.Done()
				_ = c.Handler(ctx, &mq.Message{
					Header: map[string]string{"msgID": item.MsgId, "topic": item.Topic},
					Body:   item.Body,
				})
			}, func(err error) {
				klog.WithContext(ctx).Errorf("[Client.doConsumer] err:%v", err)
			})
		} else {
			c.wg.Add(1)
			func() {
				defer c.wg.Done()
				_ = c.Handler(ctx, &mq.Message{
					Header: map[string]string{"msgID": item.MsgId, "topic": item.Topic},
					Body:   item.Body,
				})
			}()

		}
	}
	return consumer.ConsumeSuccess, nil
}

// Publish
// 	@Description
// 	@Receiver c
//	@Param ctx
//	@Param topic
//	@Param msg
// 	@Return *mq.RespMessage
// 	@Return error
func (c *Client) Publish(ctx context.Context, topic string, msg *mq.Message, opts ...mq.PublishOption) (*mq.RespMessage, error) {
	result, err := c.producerIns.SendSync(ctx, &primitive.Message{
		Topic: topic,
		Body:  msg.Body,
	})

	if err != nil {
		return nil, err
	}
	if result.Status != primitive.SendOK {
		return nil, fmt.Errorf("RocketMq producer send msg error status:%v", result.Status)
	}

	return &mq.RespMessage{
		Topic: topic,
		MsgId: result.MsgID,
	}, nil
}

// RegisterHandler
// 	@Description
// 	@Receiver c
//	@Param ctx
//	@Param h
// 	@Return error
func (c *Client) RegisterHandler(ctx context.Context, h mq.HandlerFunc) error {
	c.Handler = h
	return nil
}

// Stop
// 	@Description
// 	@Receiver c
// 	@Return error
func (c *Client) Stop() error {
	c.wg.Wait()
	return nil
}

func (c *Client) GracefulStop() error {
	c.wg.Wait()
	return nil
}

// newProducer 注册rocketmq生产者
func (c *Client) newProducer(ctx context.Context) *Client {
	addr, err := primitive.NewNamesrvAddr(c.cfg.Brokers...)
	if err != nil {
		klog.WithContext(ctx).Panic(err.Error())
	}

	c.producerIns, err = rocketmq.NewProducer(
		producer.WithNameServer(addr),
		producer.WithRetry(c.cfg.Retry),
		producer.WithGroupName(c.cfg.GroupName),
	)

	if err != nil {
		klog.WithContext(ctx).Panic(err.Error())
	}

	err = c.producerIns.Start()
	if err != nil {
		klog.WithContext(ctx).Panic(err.Error())
	}

	return c
}

// newConsumer 注册rocketmq消费者
func (c *Client) newConsumer(ctx context.Context) *Client {
	addr, err := primitive.NewNamesrvAddr(c.cfg.Brokers...)
	if err != nil {
		klog.WithContext(ctx).Panic(err.Error())
	}

	c.consumerIns, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer(addr),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(c.cfg.GroupName),
	)

	if err != nil {
		klog.WithContext(ctx).Panic(err.Error())
	}

	return c
}
