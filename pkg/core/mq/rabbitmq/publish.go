package rabbitmq

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"

	"github.com/streadway/amqp"
)

// doPublish
// 	@Description 向rabbitmq 池 发送消息
// 	@Receiver p Client
//	@Param exchange  交换机名字
//	@Param exchangeType 交换机类型
//	@Param body 消息体
//	@Param reliable
// 	@Return error 错误
func (p *Client) doPublish(ctx context.Context, exchange, exchangeType string, body []byte, options mq.RabbitPublishOptions) error {

	poolConn, aqConn, err := p.getConnect(ctx)
	if err != nil {
		return err
	}
	channel, err := aqConn.Channel()
	if err != nil {
		p.logger.WithContext(ctx).Errorf("Channel: %s", err)
		_ = p.getPool().Put(poolConn, true)
		return fmt.Errorf("channel: %s", err)
	}
	defer channel.Close()

	// 检查是否 检查是否需要 declare
	if options.IsCheckExchange {
		if err := channel.ExchangeDeclare(
			exchange,     // name
			exchangeType, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // noWait
			nil,          // arguments
		); err != nil {
			p.logger.WithContext(ctx).Errorf("Exchange Declare: %s", err)
			_ = p.getPool().Put(poolConn, true)
			return fmt.Errorf("exchange Declare：%s", err)
		}
	}
	if options.Reliable {
		if err := channel.Confirm(false); err != nil {
			p.logger.WithContext(ctx).Errorf("Channel could not be put into confirm mode: %s", err)
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}
		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer p.confirmOne(confirms)
	}

	if err = channel.Publish(
		exchange, // publish to an exchange
		"",       // routing to 0 or more queues
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		_ = p.getPool().Put(poolConn, true)
		p.logger.WithContext(ctx).Errorf("Exchange Publish: %s", err)
		return fmt.Errorf("exchange Publish: %s", err)
	}
	_ = p.getPool().Put(poolConn, false)
	return nil
}

//人们通常会保留一个发布通道、一个序列号和一组未确认的序列号并循环直到发布通道关闭。
func (p *Client) confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
		p.logger.Debugf("confirmed delivery with delivery tag: %d acked: %v", confirmed.DeliveryTag, confirmed.Ack)
	} else {
		p.logger.Errorf("failed delivery of delivery tag: %d acked: %v", confirmed.DeliveryTag, confirmed.Ack)
	}
}
