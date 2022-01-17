package main

import (
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	"github.com/LuoHongLiang0921/kuaigo/kserver"
)

func main()  {
	klog.Info()
	engine := kserver.NewEngine()
	engine.SetServerName("hahaha")
	engine.StartUp()
	engine.Run("")
}