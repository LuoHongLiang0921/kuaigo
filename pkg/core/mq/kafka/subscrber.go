package kafka

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kgo"
	"time"

	"github.com/Shopify/sarama"
)

// Subscribe
// 	@Description 订阅 主题 topic 消息，目前只支持消费者组。
// 	@Receiver c Client
//	@Param ctx 上下文
//	@Param topic 主题
// 	@Return <-chan *mq.Message
// 	@Return error
func (c *Client) Subscribe(ctx context.Context, topic string) (<-chan *mq.Message, error) {
	output := make(chan *mq.Message)

	consumeClosed, err := c.consumeMessage(ctx, topic, output)
	if err != nil {
		return nil, err
	}
	kgo.Go(func() {
		c.handleReconnects(ctx, topic, output, consumeClosed)
		close(output)
	})
	return output, nil
}

func (c *Client) handleReconnects(ctx context.Context, topic string, output chan *mq.Message, closed chan struct{}) {
	for {
		if closed != nil {
			<-closed
			c.logger.Debug("consumeMessage stopped")
		}
		select {
		case <-c.closing:
			c.logger.Debug("subscriber closed,no reconnect needed")
			return
		case <-ctx.Done():
			c.logger.Debug("ctx  cancelled,no reconnect needed")
			return
		default:
			c.logger.Debug("not closing,reconnecting")
		}
		var err error
		closed, err = c.consumeMessage(ctx, topic, output)
		if err != nil {
			c.logger.Warnf("cannot reconnect err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func (c *Client) consumeMessage(ctx context.Context, topic string, output chan *mq.Message, opts ...mq.ConsumerOption) (chan struct{}, error) {
	ctx, cancel := context.WithCancel(ctx)
	// 检测
	kgo.Go(func() {
		select {
		case <-c.closing:
			c.logger.Debug("closing subscriber,cancelling consumeMessages")
			cancel()
		case <-ctx.Done():
		}
	})
	// todo: 实现不是consumer group
	consumeMessageClosed := c.consumeGroupMessage(ctx, topic, output, opts...)
	//
	kgo.Go(func() {
		<-consumeMessageClosed
		if err := c.consumerGroup.Close(); err != nil {
			c.logger.Warnf("close consumer group failed,err:%v", err)
		}
	})
	return consumeMessageClosed, nil
}

func (c *Client) consumeGroupMessage(ctx context.Context, topic string, output chan *mq.Message, opts ...mq.ConsumerOption) chan struct{} {
	closed := make(chan struct{})
	handler := &consumerGroupHandler{
		ctx:     ctx,
		closing: c.closing,
		messageHandler: messageHandler{
			output:  output,
			logger:  c.logger,
			closing: c.closing,
		},
		wg:     c.wg,
		logger: c.logger,
	}
	kgo.Go(func() {
		err := c.consumerGroup.Consume(ctx, []string{topic}, handler)
		if err != nil {
			if err == sarama.ErrUnknown {
				c.logger.Warnf("received unknown err:%v", err)
			} else {
				c.logger.Warnf("group consumer err:%v", err)
			}
		}
		c.logger.Debugf("group consume done")
		close(closed)
	})
	return closed
}
