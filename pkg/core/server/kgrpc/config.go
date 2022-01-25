package kgrpc

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/flag"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"google.golang.org/grpc"
)

type Config struct {
	Host       string
	Port       int
	Deployment string

	Network                   string
	DisableMetric             bool
	SlowQueryThresholdInMilli int64
	ServiceAddress            string
	serverOptions             []grpc.ServerOption
	streamInterceptors        []grpc.StreamServerInterceptor
	unaryInterceptors         []grpc.UnaryServerInterceptor

	logger *klog.Logger
}

// RawConfig
// 	@Description 实例化config
//	@Param ctx 上下文
//	@Param key key
// 	@Return *Config 配置
func RawConfig(ctx context.Context, key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.WithContext(ctx).Panicf("grpc server parse config panic %v", err)
	}
	return config
}

// Build
// 	@Description
// 	@Receiver c
//	@Param ctx
// 	@Return *Server
func (c *Config) Build(ctx context.Context) *Server {
	return NewServer(ctx, c)
}

// DefaultConfig
// 	@Description 设置默认配置
// 	@Return *Config
func DefaultConfig() *Config {
	return &Config{
		Network:                   "tcp4",
		Host:                      flag.String("host"),
		Port:                      9092,
		Deployment:                constant.DefaultDeployment,
		DisableMetric:             false,
		SlowQueryThresholdInMilli: 500,
		logger:                    klog.KuaigoLogger.With(klog.FieldMod("server.grpc")),
		serverOptions:             []grpc.ServerOption{},
		streamInterceptors:        []grpc.StreamServerInterceptor{},
		unaryInterceptors:         []grpc.UnaryServerInterceptor{},
	}
}

// WithLogger
// 	@Description 设置日志
// 	@Receiver config Config
//	@Param logger 日志
// 	@Return *Config
func (config *Config) WithLogger(logger *klog.Logger) *Config {
	config.logger = logger
	return config
}

// Address
// 	@Description
// 	@Receiver config
// 	@Return string
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
