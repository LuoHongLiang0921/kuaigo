// @Description Thrift请求封装

package kthrift

import (
	"context"
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/container/pool"
	"github.com/apache/thrift/lib/go/thrift"
	"sync"
	"time"
)

var (
	ErrNoServiceName2Addr = errors.New("no Service Name to Addr mapping")
	ErrNoConnection       = errors.New("no connection exists")
)

var defaultPools = new(Pools)

type Pool struct {
	*pool.ClientPool
	ServiceName    string
	Address        string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
}

type Pools struct {
	sync.RWMutex
	once sync.Once
	// all 服务发现
	all map[string][]*Pool
	// nameConfig 服务名配置，服务名：连接地址
	nameConfig map[string]string
	next       int
}

// GetPools
// 	@Description  获取pools
// 	@Return *Pools
func GetPools() *Pools {
	return defaultPools
}

// Connect
// 	@Description 获取 thrift 连接池
//	@Param serviceName 服务名
//	@Param address 服务地址
//	@Param options 服务配置项
// 	@Return *Pool 连接
// 	@Return error 错误
func Connect(serviceName string, address string, options ...Option) (*Pool, error) {
	defaultPools.once.Do(func() {
		if defaultPools.nameConfig == nil {
			defaultPools.nameConfig = make(map[string]string)
		}
		if defaultPools.all == nil {
			defaultPools.all = make(map[string][]*Pool)
		}
	})
	return defaultPools.connect(serviceName, address, options...)
}

func (ps *Pools) connect(serviceName string, address string, options ...Option) (*Pool, error) {
	ps.RLock()
	if _, has := ps.all[address]; has {
		ps.RUnlock()
		return ps.Get(serviceName)
	}
	ps.RUnlock()
	p := newDefaultPool()
	for i := range options {
		options[i](p)
	}
	ps.Lock()
	p.Address = address
	ps.all[address] = []*Pool{p}
	ps.nameConfig[serviceName] = address
	ps.Unlock()
	return p, nil
}

// Get
// 	@Description 根据服务名获取 thrift 池
// 	@Receiver ps  Pools
//	@Param serviceName 服务名
// 	@Return *Pool 服务名对应的 thrift 池
// 	@Return error 错误
func (ps *Pools) Get(serviceName string) (*Pool, error) {
	ps.RLock()
	defer ps.RUnlock()
	addr, ok := ps.nameConfig[serviceName]
	if !ok {
		return nil, ErrNoServiceName2Addr
	}
	pools, ok := ps.all[addr]
	if !ok || len(pools) <= 0 {
		return nil, ErrNoConnection
	}
	// 负载
	for index := 0; index < len(pools); index++ {
		p := pools[ps.next]
		ps.next = (ps.next + 1) % len(pools)
		// todo: 负载均衡策略，轮训，需要实现
		//if _, err := p.Dial(); err != nil {
		//	continue
		//}
		return p, nil
	}
	return nil, ErrNoConnection
}

// Remove
// 	@Description 移除连接并释放资源
// 	@Receiver ps Pools
//	@Param serviceName 服务名
func (ps *Pools) Remove(serviceName string) {
	var addr string
	ps.RLock()
	addr = ps.nameConfig[serviceName]
	ps.RUnlock()
	ps.Lock()
	delete(ps.nameConfig, serviceName)
	pools, ok := ps.all[addr]
	if !ok {
		ps.Unlock()
		return
	}
	delete(ps.all, addr)
	ps.Unlock()
	for _, pool := range pools {
		pool.Release()
	}
}

func newDefaultPool() *Pool {
	return &Pool{
		ConnectTimeout: time.Millisecond * 100,
		ReadTimeout:    time.Millisecond * 1000,
		ClientPool: &pool.ClientPool{
			IdleTimeout: time.Minute * 30,
			MaxIdle:     256,
			Wait:        true,
		},
	}
}

// NewTStandardClient
// 	@Description 获取thrift client 信息
// 	@Receiver c Pool
// 	@Return *thrift.TSocket socket 信息
// 	@Return *thrift.TStandardClient
// 	@Return error 错误
func (c *Pool) NewTStandardClient() (*thrift.TSocket, *thrift.TStandardClient, error) {
	tSocket, err := thrift.NewTSocketTimeout(c.Address, c.ConnectTimeout)
	if err != nil {
		return nil, nil, err
	}

	if err = tSocket.Open(); err != nil {
		return nil, nil, err
	}

	if err = tSocket.SetTimeout(c.ReadTimeout); err != nil {
		return nil, nil, err
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport, err := transportFactory.GetTransport(tSocket)
	if err != nil {
		return nil, nil, err
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	protocol := protocolFactory.GetProtocol(transport)

	return tSocket, thrift.NewTStandardClient(protocol, protocol), nil
}

// WithRetry
// 	@Description 重试帮助函数，重试2次，第一次从连接池中获取，如失败，如重新第二次重建一个连接
//	@Param f 执行函数
//	@Param delay 重试间隔时间
// 	@Return error 错误
func WithRetry(f func(optionalForce ...bool) error, delay ...time.Duration) error {
	var err error
	for i := 0; i < 2; i++ {
		if i == 0 {
			err = f()
		} else {
			err = f(true)
		}

		if err == nil {
			break
		}

		if len(delay) != 0 && delay[0] != 0 {
			<-time.After(delay[0])
		}
	}
	return err
}

// CloseService
// @Description 关闭 thrift 服务
// @return error 错误
func CloseService(ctx context.Context) error {
	for k := range defaultPools.nameConfig {
		defaultPools.Remove(k)
		//for _, p := range ps {
		//	err := p.Release()
		//	if err != nil {
		//		xlog.InfoMsgf(ctx, "close thrift service %v,err:%v", k, err)
		//	}
		//}

	}
	return nil
}
