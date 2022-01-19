package kgin

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//ModName ..
const ModName = "server.gin"

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
		config.logger.Panic(context.TODO(), "http server parse config panic", klog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), klog.FieldErr(err), klog.FieldKey(key), klog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *klog.Logger) *Config {
	config.logger = logger
	return config
}

// WithHost ...
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort ...
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

func (config *Config) getContext() context.Context {
	return context.TODO()
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	server := newServer(config)
	server.engine.Use(recoverMiddleware(config.getContext(), config.logger, config.SlowQueryThresholdInMilli))

	if !config.DisableMetric {
		server.engine.Use(metricServerInterceptor())
	}

	if !config.DisableTrace {
		server.engine.Use(traceServerInterceptor())
	}
	return server
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
