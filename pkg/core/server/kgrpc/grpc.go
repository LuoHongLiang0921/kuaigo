package kgrpc

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/knet"
	"net"

	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	*grpc.Server
	listener net.Listener
	*Config
}

// NewServer
// 	@Description
//	@Param config
// 	@Return *Server
func NewServer(ctx context.Context, config *Config) *Server {
	var streamInterceptors = append(
		[]grpc.StreamServerInterceptor{defaultStreamServerInterceptor(ctx, config.logger, config.SlowQueryThresholdInMilli)},
		config.streamInterceptors...,
	)

	var unaryInterceptors = append(
		[]grpc.UnaryServerInterceptor{defaultUnaryServerInterceptor(ctx, config.logger, config.SlowQueryThresholdInMilli)},
		config.unaryInterceptors...,
	)

	config.serverOptions = append(config.serverOptions,
		grpc.StreamInterceptor(StreamInterceptorChain(streamInterceptors...)),
		grpc.UnaryInterceptor(UnaryInterceptorChain(unaryInterceptors...)),
	)

	newServer := grpc.NewServer(config.serverOptions...)
	listener, err := net.Listen(config.Network, config.Address())
	if err != nil {
		config.logger.Panic("new grpc server err", klog.FieldErrKind(ecode.ErrKindListenErr), klog.FieldErr(err))
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	fmt.Println(kcolor.Green("RPC Server run at:"))
	fmt.Printf("-  Local:   localhost:%d \r\n", config.Port)
	fmt.Printf("-  Network: %s:%d \r\n", knet.LocalIP(), config.Port)
	return &Server{
		Server:   newServer,
		listener: listener,
		Config:   config,
	}
}

func (s *Server) Serve() error {
	return s.Server.Serve(s.listener)
}

func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

func (s *Server) GracefulStop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

func (s *Server) Info() *server.ServiceInfo {
	serviceAddress := s.listener.Addr().String()
	if s.Config.ServiceAddress != "" {
		serviceAddress = s.Config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(serviceAddress),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}
