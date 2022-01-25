package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"net/http"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	ErrNoServiceName2Addr  = errors.New("no Service Name to Addr mapping")
	ErrNoConnection        = errors.New("no connection exists")
	ErrUnhealthyConnection = errors.New("unhealthy connection")
	errNoPeerPoolEntry     = errors.New("no peerPool entry")
	errNoPeerPool          = errors.New("no peerPool pool, could not connect")
	errConfigFormat        = errors.New("http config must has http:// and grpc must not, `,` Split it")

	idleDuration      = 5 * time.Second
	MaxSize           = 4 << 30
	grpcPoolSize      = 5
	ErrType           = "grpc-pool"
	defaultHealthAddr = "/health"

	pi *Pools
)

// Pool is used to manage the grpc client connection(s) for communicating with other
// task instances.  Right now it just holds one of them.
type Pool struct {
	sync.RWMutex
	// A "pool" now consists of one connection.  gRPC uses HTTP2 transport to combine
	// messages in the same TCP stream.
	conn       *grpc.ClientConn
	Addr       string
	ConfigAddr string
	lastEcho   time.Time
	ticker     *time.Ticker
	timeUsed   time.Time
	timeInit   time.Time
}

// Pools is pools.
type Pools struct {
	sync.RWMutex
	all      map[string][]*Pool
	next     int
	nameAddr map[string]string // service name => service grpc addr
	ctx      context.Context
}

// init
//  @Description 初始化
func init() {
	pi = new(Pools)
	pi.all = make(map[string][]*Pool)
	pi.nameAddr = make(map[string]string)
}

// Gets
//  @Description: 获取连接池信息
//  @Return *Pools
func Gets() *Pools {
	return pi
}

// Get
//  @Description: 通过服务名称从连接池取出一个连接，负载策略为轮训
//  @Receiver p Pools
//  @Param serviceName 服务名
//  @Return *Pool 服务名连接
//  @Return error 错误
func (p *Pools) Get(serviceName string) (*Pool, error) {
	p.RLock()
	defer p.RUnlock()
	addr, ok := p.nameAddr[serviceName]
	if !ok {
		return nil, ErrNoServiceName2Addr
	}
	pools, ok := p.all[addr]
	if !ok || len(pools) <= 0 {
		return nil, ErrNoConnection
	}
	for index := 0; index < len(pools); index++ {
		pool := pools[p.next]
		p.next = (p.next + 1) % len(pools)
		if !pool.IsHealthy() {
			continue
		}
		pool.timeUsed = time.Now()
		return pool, nil
	}
	return nil, ErrNoConnection
}

func (p *Pools) WithContext(ctx context.Context) *Pools {
	p.ctx = ctx
	return p
}

// getContext
//  @Description: 获取上下文
//  @Receiver p Pools
//  @Return context.Context 上下文
func (p *Pools) getContext() context.Context {
	if p.ctx == nil {
		return context.TODO()
	}
	return p.ctx
}

// Connect
//  @Description: 创建链接
//  @Receiver p
//  @Param serviceName 服务名称
//  @Param configAddr rpc地址 http开头 多个通过,分隔
//  @Return *Pool
//  @Return error
func (p *Pools) Connect(serviceName, configAddr string) (*Pool, error) {
	addr, _, err := addrCheck(configAddr)
	if err != nil {
		return nil, err
	}

	p.RLock()
	if _, has := p.all[addr]; has {
		p.RUnlock()
		return p.Get(serviceName)
	}
	p.RUnlock()
	var pools []*Pool
	for index := 0; index < grpcPoolSize; index++ {
		pool, err := createClient(addr, configAddr)
		if err != nil {
			//log.Error(ErrType).Msgf("Unable to connect to host: %s, err: %v", addr, err)
			klog.WithContext(p.getContext()).Error("GetWithUnmarshal.Decode",
				klog.String("name", "grpc-pool"),
				klog.String("msg", fmt.Sprintf("Unable to connect to host: %s, err: %v", addr, err)),
				klog.FieldErr(err))
			return nil, fmt.Errorf("connect addr %v err: %v", addr, err)
		}
		pool.ConfigAddr = configAddr
		pools = append(pools, pool)
	}
	if len(pools) > 0 {
		p.Lock()
		p.all[addr] = pools
		p.nameAddr[serviceName] = addr
		p.Unlock()
	}
	return p.Get(serviceName)
}

// createClient
//  @Description: 创建新的连接池
//  @Param addr
//  @Param configAddr
//  @Return *Pool
//  @Return error
func createClient(addr, configAddr string) (*Pool, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(MaxSize),
			grpc.MaxCallSendMsgSize(MaxSize)),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	now := time.Now()
	pl := &Pool{
		conn:       conn,
		Addr:       addr,
		ConfigAddr: configAddr,
		timeInit:   now,
		timeUsed:   now,
	}

	if !pl.health() {
		return nil, errors.New("http health check failed")
	}
	pl.UpdateHealthStatus()
	go pl.MonitorHealth()
	return pl, nil
}

// Remove
//  @Description: 关闭指定serviceName的连接池
//  @Receiver p
//  @Param serviceName
//  @Param addr
func (p *Pools) Remove(serviceName, addr string) {
	p.Lock()
	delete(p.nameAddr, serviceName)
	addr, _, err := addrCheck(addr)
	if err != nil {
		return
	}
	pools, ok := p.all[addr]
	if !ok {
		p.Unlock()
		return
	}
	delete(p.all, addr)
	p.Unlock()
	for _, pool := range pools {
		pool.close()
	}
	return
}

// addrCheck
//  @Description: 检查配置地址正确性
//  @Param addr 地址 格式 http://XX.XX.XX,XXX.XXX.XXX:XXXX
//  @Return string grpc地址
//  @Return string http地址
//  @Return error 错误信息
func addrCheck(addr string) (string, string, error) {
	ss := strings.Split(addr, ",")
	if len(ss) > 2 {
		return "", "", errConfigFormat
	}
	hasGrpc, hasHttp := false, false
	var grpcAddr, httpAddr string
	for _, s := range ss {
		if strings.HasPrefix(s, "http://") {
			hasHttp = true
			httpAddr = s
		} else {
			hasGrpc = true
			grpcAddr = s
		}
	}
	//如果GRPC和HTTP配置同事不存在
	if !hasGrpc && !hasHttp {
		return "", "", errConfigFormat
	}
	return grpcAddr, httpAddr, nil
}

// close
//  @Description: 关闭连接池
//  @Receiver p
func (p *Pool) close() {
	p.ticker.Stop()
	p.conn.Close()
}

// UpdateHealthStatus
//  @Description: 更新健康状态
//  @Receiver p
func (p *Pool) UpdateHealthStatus() {
	p.Lock()
	p.lastEcho = time.Now()
	p.Unlock()
}

// MonitorHealth
//  @Description: 监测健康状态
//  @Receiver p
func (p *Pool) MonitorHealth() {
	p.ticker = time.NewTicker(idleDuration)
	for range p.ticker.C {
		if p.health() {
			p.UpdateHealthStatus()
		}
	}
}

// health
//  @Description: 健康检查
//  @Receiver p
//  @Return bool
func (p *Pool) health() bool {
	var netClient = &http.Client{
		Timeout: 1000 * time.Millisecond,
	}
	_, httpAddr, err := addrCheck(p.ConfigAddr)
	if httpAddr == "" {
		//如果没有配置http健康检查则志杰返回 todo
		return true
	} else if err != nil {
		return false
	}
	resp, err := netClient.Get(httpAddr + defaultHealthAddr)
	if err == nil && resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

// IsHealthy
//  @Description: 健康状态是否良好
//  @Receiver p
//  @Return bool
func (p *Pool) IsHealthy() bool {
	p.RLock()
	defer p.RUnlock()
	return time.Since(p.lastEcho) < 2*idleDuration
}

// Conn returns the connection to use from the pool of connections.
func (p *Pool) Conn() *grpc.ClientConn {
	p.RLock()
	defer p.RUnlock()
	return p.conn
}

// Get returns the connection to use from the pool of connections.
func (p *Pool) Get() *grpc.ClientConn {
	p.RLock()
	defer p.RUnlock()
	return p.conn
}

// IsClose
//  @Description: 连接池是否关闭
//  @Receiver pool
//  @Return bool
func (p *Pools) IsClose(serviceName string) (bool, error) {
	p.RLock()
	defer p.RUnlock()
	addr, ok := p.nameAddr[serviceName]
	if !ok {
		return true, ErrNoServiceName2Addr
	}
	pools, ok := p.all[addr]
	if !ok || len(pools) <= 0 {
		return true, ErrNoConnection
	}
	return false, nil
}

// TimeInit
//  @Description: 获取连接创建时间
//  @Receiver client
//  @Return time.Time
func (p *Pool) TimeInit() time.Time {
	return p.timeInit
}

// TimeUsed
//  @Description: 获取连接上一次使用时间
//  @Receiver client
//  @Return time.Time
func (p *Pool) TimeUsed() time.Time {
	return p.timeUsed
}
