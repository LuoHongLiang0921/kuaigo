package mq

import (
	"context"
	"time"
)

// MessageQueer 消息队列接口
type MessageQueer interface {
	Consumer
	Publisher
	Handler
	Stoper
	GracefulStoper
}

// RespMessage 响应消息
type RespMessage struct {
	Topic     string    `json:"topic"`
	MsgId     string    `json:"msg_id"`
	Offset    int64     `json:"offset,omitempty"`
	Partition int32     `json:"partition,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// Consumer 消费接口
type Consumer interface {
	Consume(ctx context.Context, opts ...ConsumerOption) error
}

// Subscriber 订阅接口
type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
}

//Publisher 发布消息
type Publisher interface {
	// Publish
	// 	@Description
	//	@Param ctx 上下文
	//	@Param target kafka 为topic,rabbit mq 为 queue
	//	@Param msg 消息内容
	//	@Param opt 发布配置项
	// 	@Return *RespMessage 发布消息后的响应内容
	// 	@Return error 错误
	Publish(ctx context.Context, target string, msg *Message, opt ...PublishOption) (*RespMessage, error)
}

// Handler 逻辑处理
type Handler interface {
	RegisterHandler(ctx context.Context, h HandlerFunc) error
}

// Stoper  停止
type Stoper interface {
	Stop() error
}

// GracefulStoper 优雅停止
type GracefulStoper interface {
	GracefulStop() error
}

// HandlerFunc  业务逻辑处理函数
type HandlerFunc func(ctx context.Context, msg *Message) error
