package ginserver

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/controller"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/middleware"
	"net/http"

)

type GinServer struct {
	serverName          string
	controllerGroupName string
	*Server
	Middleware middleware.IMiddleware
}

// NewServer
//  @Description 实例化路由
//  @Param name 服务名称
//  @Param controllerGroupName controller分组名称
//  @Param protocol 服务协议，配置文件中别名
//  @Return GinServer 本实例
func NewServer(serverName, controllerGroupName, protocol string) *GinServer {
	r := new(GinServer)
	r.serverName = serverName
	r.controllerGroupName = controllerGroupName
	r.Server = RawConfig("listen." + protocol).Build()
	return r
}

// WithMiddleware
//  @Description 增加中间件注入
//  @Receiver s GinServer
//  @Param m 中间件实例
//  @Return GinServer 本实例
func (s *GinServer) WithMiddleware(m middleware.IMiddleware) *GinServer {
	s.Middleware = m
	return s
}

// Build
//  @Description 构建服务配置
//  @Receiver s GinServer实例
//  @Param fns 初始化函数集合
func (s *GinServer) Build() *GinServer {
	for _, c := range controller.GetControllers(s.controllerGroupName) {
		c.RegisterController(s, s.Middleware)
	}
	return s
}

func (s *Server) Use(middleware ...HandlerFunc) IRoutes {
	return s.engine.Use(middleware...)
}

func (s *Server) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.Handle(httpMethod, relativePath, handlers...)
}

func (s *Server) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.Any(relativePath, handlers...)
}

func (s *Server) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.GET(relativePath, handlers...)
}

func (s *Server) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.POST(relativePath, handlers...)
}

func (s *Server) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.DELETE(relativePath, handlers...)
}

func (s *Server) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.PATCH(relativePath, handlers...)
}

func (s *Server) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.PUT(relativePath, handlers...)
}

func (s *Server) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.PUT(relativePath, handlers...)
}

func (s *Server) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return s.engine.HEAD(relativePath, handlers...)
}

func (s *Server) StaticFile(relativePath, filepath string) IRoutes {
	return s.engine.StaticFile(relativePath, filepath)
}

func (s *Server) Static(relativePath, root string) IRoutes {
	return s.engine.Static(relativePath, root)
}

func (s *Server) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	return s.engine.StaticFS(relativePath, fs)
}

func (s *Server) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return s.engine.Group(relativePath, handlers...)
}
