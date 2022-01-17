package core

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kutils"
	"sync"
)

type App struct {
	smu          *sync.RWMutex
	initOnce     sync.Once
	startupOnce  sync.Once
	stopOnce     sync.Once
	modsOnce     sync.Once
	//ServiceName 上层服务名
	ServiceName string
}

func (app *App) StartUp() error {
	// 启动可以做的事情
	var err error
	app.startupOnce.Do(func() {
		err = kutils.SerialUntilError(
			app.printBanner,
			//app.initLog,
		)()
	})
	return err
}

func (app *App)printBanner() error {
	//if app.HideBanner {
	//	return nil
	//}
	const banner = `Welcome to my frame, starting application ...`
	//fmt.Println(xcolor.Green(banner))
	fmt.Println(fmt.Sprintf("\x1b[32m%s\x1b[0m", banner))
	return nil
}
