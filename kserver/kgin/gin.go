package kgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	TContext    = gin.Context
	Engine      = gin.Engine
	HandlerFunc = gin.HandlerFunc
	RouterGroup = gin.RouterGroup
	IRoutes     = gin.IRoutes
)

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
