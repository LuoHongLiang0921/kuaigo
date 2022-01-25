package rabbitmq

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kgo"
	"time"

	"github.com/streadway/amqp"
)

// Subscribe
// 	@Description 订阅 queue 队列消息
// 	@Receiver p Client
//	@Param ctx 上下文
//	@Param topic 队列名
// 	@Return <-chan 消息
// 	@Return error
func (p *Client) Subscribe(ctx context.Context, queue string) (<-chan *mq.Message, error) {
	output := make(chan *mq.Message)
	consumeClosed := make(chan struct{})

	var err error
	consumeClosed, err = p.consumeMessage(ctx, queue, output)
	if err != nil {
		p.logger.Warnf("consumeMessage err:%v", err)
		return nil, err
	}

	kgo.Go(func() {
		<-consumeClosed
		p.handleReconnects(ctx, queue, output, consumeClosed)
		close(output)
	})
	return output, nil
}

func (p *Client) consumeMessage(ctx context.Context, queue string, output chan *mq.Message) (chan struct{}, error) {
	closed := make(chan struct{})
	// get connect
	poolConn, aqConn, err := p.getConnect(ctx)
	if err != nil {
		close(closed)
		return closed, err
	}
	// get channel
	aqChannel, err := aqConn.Channel()
	if err != nil {
		close(closed)
		_ = p.getPool().Put(poolConn, true)
		return closed, nil
	}
	// createConsumer
	delivery, err := p.createConsumer(ctx, queue, aqChannel)
	if err != nil {
		return closed, err
	}
	// notify close error signal
	notifyClosing := make(chan struct{})
	p.processNotifyClose(ctx, aqConn, notifyClosing)
	kgo.Go(func() {
		for {
			select {
			case d, ok := <-delivery:
				if !ok {
					p.logger.Debug("rabbit mq delivery closing,prepare reconnecting")
					close(closed)
					return
				}
				msg := mq.NewMessage(d.Body)
				output <- msg
				select {
				case <-msg.Acked():
					_ = d.Ack(false)
				case <-msg.NAcked():
					_ = d.Nack(false, p.cfg.ConsumeConfig.IsNackRequeue)
				}
			case <-notifyClosing:
				p.logger.Debug("rabbit mq closing,prepare reconnecting")
				close(closed)
				return
			case <-p.closing:
				p.logger.Debug("consumeMessage stop")
				return
			}
		}
	})
	return closed, nil
}

func (p *Client) createConsumer(ctx context.Context, queue string, aqChannel *amqp.Channel) (<-chan amqp.Delivery, error) {
	return aqChannel.Consume(
		queue,
		p.cfg.ConsumeConfig.ConsumerTag,
		p.cfg.ConsumeConfig.AutoAck,
		p.cfg.ConsumeConfig.Exclusive,
		p.cfg.ConsumeConfig.NoLocal,
		p.cfg.ConsumeConfig.NoWait,
		p.cfg.ConsumeConfig.Args,
	)
}

func (p *Client) processNotifyClose(ctx context.Context, aqConn *amqp.Connection, closing chan struct{}) {
	kgo.Go(func() {
		notifyClosed := make(chan *amqp.Error)
		connClosed := aqConn.NotifyClose(notifyClosed)
		err := <-connClosed
		p.logger.Warnf("receiver close notification from rabbit, err %v", err)
		close(closing)
		p.logger.Debug("receiver close notification from rabbit--->")
	})

}

func (p *Client) handleReconnects(ctx context.Context, queue string, output chan *mq.Message, closed chan struct{}) {
	for {
		if closed != nil {
			<-closed
			p.logger.Debug("consumeMessage stopped")
		}
		select {
		case <-p.closing:
			p.logger.Debug("subscriber closed,no reconnect needed")
			return
		case <-ctx.Done():
			p.logger.Debug("ctx  cancelled,no reconnect needed")
			return
		default:
			p.logger.Debug("not closing,reconnecting")
		}
		var err error
		closed, err = p.consumeMessage(ctx, queue, output)
		if err != nil {
			p.logger.Warnf("cannot reconnect err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
}
