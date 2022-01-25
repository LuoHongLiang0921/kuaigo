// @Description 系统环境变量

package pkg

import (
	"crypto/md5"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"os"
)

var (
	appLogDir            string
	appMode              string
	appRegion            string
	appZone              string
	appHost              string
	appInstance          string
	apolloServiceAddress string
)

func InitEnv() {
	appLogDir = os.Getenv(constant.EnvAppLogDir)
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appHost = os.Getenv(constant.EnvAppHost)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", GetHostName(), GetAppID()))))
	}
}

func AppLogDir() string {
	return appLogDir
}

func SetAppLogDir(logDir string) {
	appLogDir = logDir
}

func GetAppMode() string {
	return appMode
}

func SetAppMode(mode string) {
	appMode = mode
}

func GetAppRegion() string {
	return appRegion
}

func SetAppRegion(region string) {
	appRegion = region
}

func GetAppZone() string {
	return appZone
}

func SetAppZone(zone string) {
	appZone = zone
}

func GetAppHost() string {
	return appHost
}

func SetAppHost(host string) {
	appHost = host
}

func AppInstance() string {
	return appInstance
}
