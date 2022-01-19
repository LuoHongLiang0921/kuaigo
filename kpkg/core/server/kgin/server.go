package kgin

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/server"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
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

func newServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		config.logger.Panic(config.getContext(), "new xgin server err", klog.FieldErrKind(ecode.ErrKindListenErr), klog.FieldErr(err))
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

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	// s.Gin.StdLogger = xlog.TabbyLogger.StdLog()
	//for _, route := range s.engine.Routes() {
	//	s.config.logger.Info(context.TODO(), "add route", xlog.FieldMethod(route.Method), xlog.String("path", route.Path))
	//}
	s.Server = &http.Server{
		Addr:    s.config.Address(),
		Handler: s.engine,
	}
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		//s.config.logger.Info(context.TODO(), "close gin", xlog.FieldAddr(s.config.Address()))
		return nil
	}

	return err
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Info returns server info, used by governor and consumer balancer
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

//GetEngine 返回 server engine
func (s *Server) SetLogger(l *klog.Logger) {
	s.config.logger = l
}
