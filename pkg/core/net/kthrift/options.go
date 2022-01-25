// @Description thrift封装

package kthrift

import "time"

// Option Option 配置函数
type Option func(c *Pool)

// WithConnectTimeout conn timeout
// WithConnectTimeout
// 	@Description
//	@Param t
// 	@Return Option
func WithConnectTimeout(t time.Duration) Option {
	return func(p *Pool) {
		p.ConnectTimeout = t
	}
}

// WithReadTimeout read timeout
// WithReadTimeout
// 	@Description
//	@Param t
// 	@Return Option
func WithReadTimeout(t time.Duration) Option {
	return func(p *Pool) {
		p.ReadTimeout = t
	}
}

// WithMaxIdle max idl
// WithMaxIdle
// 	@Description
//	@Param n
// 	@Return Option
func WithMaxIdle(n uint32) Option {
	return func(p *Pool) {
		p.MaxIdle = int(n)
	}
}

// WithIdleTimeout idl timeout
// WithIdleTimeout
// 	@Description
//	@Param t
// 	@Return Option
func WithIdleTimeout(t time.Duration) Option {
	return func(p *Pool) {
		p.IdleTimeout = t
	}
}

// WithMaxActive set max active conn in pool
// WithMaxActive
// 	@Description
//	@Param n
// 	@Return Option
func WithMaxActive(n int) Option {
	return func(p *Pool) {
		p.MaxActive = n
	}
}

// WithMaxLiveTime is conn max lifetime
// WithMaxLiveTime
// 	@Description
//	@Param t
// 	@Return Option
func WithMaxLiveTime(t time.Duration) Option {
	return func(p *Pool) {
		p.MaxLiveTime = t
	}
}
