package kgin

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"net/http"

	"net"

	"github.com/gin-gonic/gin"
)

// Server ...
type Server struct {
	engine   *Engine
	Server   *http.Server
	config   *Config
	listener net.Listener
}

// newServer
//  @Description  实例化gin server
//  @Param config 配置
//  @Return *Server
func newServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		config.logger.Panic("new xgin server err", klog.FieldErrKind(ecode.ErrKindListenErr), klog.FieldErr(err))
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	gin.SetMode(config.Mode)
	return &Server{
		engine:   gin.New(),
		config:   config,
		listener: listener,
	}
}

//Upgrade protocol to WebSocket
func (s *Server) Upgrade(ws *WebSocket) gin.IRoutes {
	return s.engine.GET(ws.Pattern, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request)
	})
}

// Serve 实现serve接口
//  @Description  启动server服务
//  @Receiver s
//  @Return error
func (s *Server) Serve() error {
	// s.Gin.StdLogger = xlog.TabbyLogger.StdLog()
	for _, route := range s.engine.Routes() {
		s.config.logger.Info("add route", klog.FieldMethod(route.Method), klog.String("path", route.Path))
	}
	s.Server = &http.Server{
		Addr:    s.config.Address(),
		Handler: s.engine,
	}
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		s.config.logger.Info("close gin", klog.FieldAddr(s.config.Address()))
		return nil
	}

	return err
}

// Stop 实现stop接口
//  @Description 立即终止gin服务器
//  @Receiver s
//  @Return error
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop 实现GracefulStop接口
//  @Description  优雅停止gin服务器
//  @Receiver s
//  @Param ctx
//  @Return error
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Info
//  @Description  初始化server信息
//  @Receiver s
//  @Return *server.ServiceInfo
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.listener.Addr().String()
	if s.config.ServiceAddress != "" {
		serviceAddr = s.config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	// info.Name = info.Name + "." + ModName
	return &info
}

//SetLogger 设置日志
func (s *Server) SetLogger(l *klog.Logger) {
	s.config.logger = l
}
