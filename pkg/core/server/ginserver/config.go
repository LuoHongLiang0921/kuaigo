package ginserver

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/knet"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//ModName ..
const ModName = "ginserver"

type (
	TContext    = gin.Context
	Engine      = gin.Engine
	HandlerFunc = gin.HandlerFunc
	RouterGroup = gin.RouterGroup
	IRoutes     = gin.IRoutes
	IRouter     = gin.IRouter
)

// Config HTTP config
type Config struct {
	Host          string
	Port          int
	Deployment    string
	Mode          string
	DisableMetric bool
	DisableTrace  bool
	// ServiceAddress service address in registry info, default to 'Host:Port'
	ServiceAddress string

	SlowQueryThresholdInMilli int64

	logger *klog.Logger
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                      "127.0.0.1",
		Port:                      9091,
		Mode:                      gin.ReleaseMode,
		SlowQueryThresholdInMilli: 500, // 500ms
		logger:                    klog.KuaigoLogger.With(klog.FieldMod(ModName)),
	}
}

// StdConfig Jupiter Standard HTTP Server config
func StdConfig(name string) *Config {
	return RawConfig("tabby.server." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil &&
		errors.Cause(err) != conf.ErrInvalidKey {
		config.logger.Panic("ginserver parse config panic", klog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), klog.FieldErr(err), klog.FieldKey(key), klog.FieldValueAny(config))
	}
	return config
}
func (config *Config) getContext() context.Context {
	return context.TODO()
}

// Build
//  @Description  构建并创建服务器实例，使用必要的拦截器对其进行初始化
//  @Receiver config
//  @Return *Server
func (config *Config) Build() *Server {
	server := newServer(config)
	server.engine.Use(recoverMiddleware(config.getContext(), config.logger, config.SlowQueryThresholdInMilli))

	if !config.DisableMetric {
		server.engine.Use(metricServerInterceptor())
	}

	fmt.Println(kcolor.Green("Web Server run at:"))
	fmt.Printf("-  Local:   http://localhost:%d/ \r\n", config.Port)
	fmt.Printf("-  Network: http://%s:%d/ \r\n", knet.LocalIP(), config.Port)
	return server
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
