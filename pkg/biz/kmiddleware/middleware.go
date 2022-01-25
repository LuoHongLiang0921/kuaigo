package kmiddleware

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/ratelimiter"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/ginserver"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/storage/redis"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kjwt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"time"

)

type Option func(m *Middleware)

// Middleware 统一中间件
type Middleware struct {
	AllowCloseSignValidate bool
	logger                 *klog.Logger
	rater                  *ratelimiter.RateLimiter
	// 解析用户id
	jwt      *kjwt.Client
	adminJwt *kjwt.Client
	//ABTest 来源
	ABSource   string
	GuestModel GuestUserModel
}

const (
	RatelimiterConfigKey = "middleware.rateLimiter"
	cacheConfigPrefix    = "caches."
)

// NewMiddleware
//  @Description 实例化中间件服务
//  @Param opts
//  @Return *Middleware
func NewMiddleware(opts ...Option) *Middleware {
	m := new(Middleware)
	for _, f := range opts {
		f(m)
	}
	return m
}

// WithIgnoreSignValidate
//  @Description: 忽略签名验证
//  @Param flag
//  @Return Option
func WithIgnoreSignValidate(flag bool) Option {
	return func(m *Middleware) {
		m.AllowCloseSignValidate = flag
	}
}

// WithLogger
//  @Description 开启日志中间件
//  @Param l
//  @Return Option
func WithLogger(l *klog.Logger) Option {
	return func(m *Middleware) {
		m.logger = l
	}
}

// WithJwt
//  @Description: 开启JWT中间件
//  @Param jwtKey jwt密钥
//  @Param expireTime 过期时间 单位s
//  @Return Option
func WithJwt(jwtKey string, expireTime int64) Option {
	return func(m *Middleware) {
		m.jwt = &kjwt.Client{
			JwtKey: jwtKey,
			Expire: time.Duration(expireTime) * time.Second,
		}
	}
}

// WithAdminJwt
//  @Description 开启后台JWT中间件
//  @Param jwt
//  @Return Option
func WithAdminJwt(jwtKey string, expireTime int64) Option {
	return func(m *Middleware) {
		m.adminJwt = &kjwt.Client{
			JwtKey: jwtKey,
			Expire: time.Duration(expireTime) * time.Second,
		}
	}
}

// WithRateLimiter
//  @Description: 开启限流中间件
//  @Param redisKey 配置文件中用于限流的rediskey
//  @Return Option
func WithRateLimiter(redisKey string) Option {
	return func(m *Middleware) {
		m.rater = ratelimiter.RawConfig(RatelimiterConfigKey).
			Build().
			WithRedis(redis.RawRedisConfig(cacheConfigPrefix + redisKey).Build())
	}
}

//// WithABTest AbSource 实例
//// 	@Description: AbSource 实例
////	@Param a
//// 	@return Option
//func WithABTest(source string) Option {
//	return func(m *Middleware) {
//		m.ABSource = source
//	}
//}

// WithGuestModel user实例
// 	@Description user实例
//	@Param userModel
// 	@return Option
func WithGuestModel(redisKey, addr string) Option {
	return func(m *Middleware) {
		m.GuestModel = NewGuestUserModel(redis.RawRedisConfig(cacheConfigPrefix+redisKey).Build(), addr)
	}
}

// AppRatelimt
//  @Description  app 限流
//  @Receiver m
//  @Param c
func (m *Middleware) AppRatelimt(c *ginserver.TContext) {
	m.ratelimt(c)
}

// AdminRatelimt
//  @Description  后台限流
//  @Receiver m
//  @Param c
func (m *Middleware) AdminRatelimt(c *ginserver.TContext) {
	m.ratelimt(c)
}

func (m *Middleware) ratelimt(c *ginserver.TContext) {

}

// Deprecated: 请使用 WithIgnoreSignValidate 替代
//  @Description  关闭header验签
//  @Param flag
//  @Return Option
func WithAllowCloseSignValidate(flag bool) Option {
	return WithIgnoreSignValidate(flag)
}
