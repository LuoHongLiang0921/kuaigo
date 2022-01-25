package kafka

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"

	"github.com/Shopify/sarama"
)

// Consumer ...
type consumerGroupHandler struct {
	ctx            context.Context
	closing        chan struct{}
	messageHandler messageHandler
	logger         *klog.Logger
	wg             sync.WaitGroup
}

// Setup
// 	@Description consumer group 分区规则，实现 ConsumerGroupHandler 接口
// 	@Receiver c Consumer
//	@Param session 消费组会话
// 	@Return error 错误
func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup
// 	@Description 会话结束阶段，实现 ConsumerGroupHandler 接口
// 	@Receiver c Consumer
//	@Param session 消费组会话
// 	@Return error 错误
func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim
// 	@Description 消费消息
// 	@Receiver c Consumer
//	@Param session 消费组会话
//	@Param claim 消息组
// 	@Return error
func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	kafkaMsgs := claim.Messages()
	for {
		select {
		case kafkaMsg, ok := <-kafkaMsgs:
			if !ok {
				c.logger.Debug("kafka messages closed,stopping consumerGroupHandler")
				return nil
			}
			c.wg.Add(1)
			err := c.messageHandler.processMessage(c.ctx, kafkaMsg, session)
			if err != nil {
				c.wg.Done()
				return err
			}
			c.wg.Done()
		case <-c.closing:
			c.logger.Debug("subscriber  closed,stopping consumerGroupHandler")
			return nil
		case <-c.ctx.Done():
			c.logger.Debug("ctx was cancelled,stopping consumerGroupHandler")
			return nil

		}
	}
}

type messageHandler struct {
	output  chan<- *mq.Message
	logger  *klog.Logger
	closing chan struct{}
}

func (h messageHandler) processMessage(ctx context.Context, kafkaMessage *sarama.ConsumerMessage, sess sarama.ConsumerGroupSession) error {
	msg := mq.NewMessage(kafkaMessage.Value)
	ctx, cancel := context.WithCancel(ctx)
	msg.SetContext(ctx)
	defer cancel()
loop:
	for {
		select {
		case h.output <- msg:
		case <-h.closing:
			return nil
		case <-ctx.Done():
			return nil
		}

		select {
		case <-msg.Acked():
			if sess != nil {
				sess.MarkMessage(kafkaMessage, "")
			}
			break loop
		case <-msg.NAcked():
			msg = msg.Copy()
			break loop
		case <-h.closing:
			return nil
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}
