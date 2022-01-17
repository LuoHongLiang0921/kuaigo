package kserver

import (
	"github.com/LuoHongLiang0921/kuaigo/core"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	core.App
	Eng *gin.Engine
}
func NewEngine() *Engine {
	eng := Engine{}
	eng.Eng = gin.New()
	return &eng
}
func (e *Engine) SetServerName(s string)  {
	e.ServiceName = s
}

func (e *Engine) Run(port string)  {
	runPort := ""
	defaultPort := "8080"
	if port != ""{
		runPort = ":"+ port
	} else {
		runPort = ":"+ defaultPort
	}

	e.Eng.Run(runPort)
}