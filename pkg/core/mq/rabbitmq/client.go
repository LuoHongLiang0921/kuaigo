package rabbitmq

import (
	"context"
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kpool"
	"sync"

	"github.com/streadway/amqp"
)

// Client ...
type Client struct {
	cfg     *Config
	Handler mq.HandlerFunc

	logger *klog.Logger
	// closing 停止从rabbitmq 接受新的消息
	closing chan struct{}
	wg      sync.WaitGroup

	mu  sync.Mutex
	err error

	poolMu sync.RWMutex
	pool   *kpool.Pool
}

func (p *Client) processMessage(ctx context.Context, mqMsg *mq.Message) error {
	p.wg.Add(1)
	err := p.Handler(ctx, mqMsg)
	if err != nil {
		p.wg.Done()
		p.logger.WithContext(ctx).Warnf("process addr %v,queue %v is err:%v", p.cfg.Address, p.cfg.ConsumeConfig.Queue, err)
		return err
	}
	p.wg.Done()
	return nil
}

// Consume
// 	@Description 消费消息，消息需要 mq.Message  中的 Ack() 或 NAck() 方法，否则会阻塞。
// 	@Receiver p Client
//	@Param ctx 上下文
//	@Param opts 消息配置项
// 	@Return error 错误
func (p *Client) Consume(ctx context.Context, opts ...mq.ConsumerOption) error {
	consumerOptions := mq.MergeConsumerOption(opts...)
	queue := p.cfg.ConsumeConfig.Queue
	if consumerOptions.Queue != "" {
		queue = consumerOptions.Queue
	}
	mqMsgs, err := p.Subscribe(ctx, queue)
	if err != nil {
		return err
	}
	for {
		select {
		case msg, ok := <-mqMsgs:
			if !ok {
				p.logger.Warnf("message consume closed")
				return nil
			}
			_ = p.processMessage(ctx, msg)
		case <-p.closing:
			return nil
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *Client) getPool() *kpool.Pool {
	p.poolMu.RLock()
	defer p.poolMu.RUnlock()
	return p.pool
}

func (p *Client) getConnect(ctx context.Context) (*kpool.PoolConn, *amqp.Connection, error) {
	poolConn, err := p.getPool().Get(nil)
	if err != nil {
		_ = p.getPool().Put(poolConn, true)
		return nil, nil, err
	}
	connection, ok := poolConn.C.(*amqp.Connection)
	if !ok {
		_ = p.getPool().Put(poolConn, true)
		return poolConn, nil, errors.New("connection isn't amqp.Connection")
	}
	return poolConn, connection, nil
}

// Publish
// 	@Description 发布消息
// 	@Receiver p Client
//	@Param ctx 上下文
//	@Param topic 主题
//	@Param msg 消息
//	@Param opts 发布配置项
// 	@Return *mq.RespMessage 发布消息的响应信息
// 	@Return error
func (p *Client) Publish(ctx context.Context, target string, msg *mq.Message, opts ...mq.PublishOption) (*mq.RespMessage, error) {
	publishOpt := mq.MergePublishOption(opts...)
	err := p.doPublish(ctx, target, publishOpt.ExchangeType, msg.Body, publishOpt.RabbitPublishOptions)
	if err != nil {
		return nil, err
	}
	return &mq.RespMessage{}, nil
}

// RegisterHandler
// 	@Description 注册业务函数
// 	@Receiver p Client
//	@Param ctx 上下文
//	@Param h 执行函数
// 	@Return error 错误
func (p *Client) RegisterHandler(ctx context.Context, h mq.HandlerFunc) error {
	p.Handler = h
	return nil
}

func (p *Client) Stop() error {
	p.notifyConsumerStop()
	return p.getPool().Close()
}

func (p *Client) GracefulStop() error {
	p.notifyConsumerStop()
	p.wg.Wait()
	return nil
}

func (p *Client) notifyConsumerStop() {
	if p.closing != nil {
		close(p.closing)
	}
}

// Close
// 	@Description 关闭队列池
// 	@Receiver p Client
// 	@Return error
func (p *Client) Close() error {
	p.getPool().Close()
	return nil
}

// Err
// 	@Description 获取池中错误
// 	@Receiver p Client
// 	@Return error 错误
func (p *Client) Err() error {
	p.mu.Lock()
	err := p.err
	p.mu.Unlock()
	return err
}
