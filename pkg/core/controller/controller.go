package controller

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/middleware"
	"sync"

	"github.com/gin-gonic/gin"

)

var controllerMap map[string][]IController
var smu sync.RWMutex

type IController interface {
	// RegisterController
	// 	@Description 注册路由，中间件
	//	@Param r 路由
	//	@Param m 中间件
	RegisterController(r gin.IRouter, m middleware.IMiddleware)
}

func init() {
	controllerMap = make(map[string][]IController)
}

// RegisterController
//  @Description 注册IController接口集合
//  @Param controllerGroupName controller分组名称
//  @Param c controller集合
func RegisterController(controllerGroupName string, c ...IController) {
	smu.Lock()
	controllers, ok := controllerMap[controllerGroupName]
	if !ok {
		controllerMap[controllerGroupName] = make([]IController, 0)
	}
	controllers = append(controllers, c...)
	controllerMap[controllerGroupName] = controllers
	smu.Unlock()
}

// GetControllers
//  @Description 获取IController接口集合
//  @Param controllerGroupName controller分组名称
//  @Return IController接口集合
func GetControllers(controllerGroupName string) []IController {
	controllers, ok := controllerMap[controllerGroupName]
	if ok {
		return controllers
	}
	return nil
}
