// @Description 系统

package pkg

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const tabbyVersion = "0.0.6"

var (
	startTime string
	goVersion string
)

// build info
var (
	appID         string
	AppName       string // 应用名字
	BuildVersion  string // build commit info
	BuildUser     string // 提交者
	BuildTag      string // tag 名
	BuildHostName string // host name
	BuildStatus   string // status
	BuildDate     string // 构建时间
	BuildBranch   string // 构建分支
	BuildGitURL   string // 仓库地址
)

func init() {
	if AppName == "" {
		AppName = os.Getenv(constant.EnvAppName)
		if AppName == "" {
			AppName = filepath.Base(os.Args[0])
		}
	}

	if BuildHostName == "" {
		name, err := os.Hostname()
		if err != nil {
			name = "unknown"
		}
		BuildHostName = name
	}

	startTime = ktime.TS.Format(time.Now())
	goVersion = runtime.Version()
	InitEnv()
}

// GetAppName gets application name.
func GetAppName() string {
	return AppName
}

//SetAppName set app anme
func SetAppName(s string) {
	AppName = s
}

//GetAppID get appID
func GetAppID() string {
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

//GetAppVersion get buildAppVersion
func GetAppVersion() string {
	return BuildTag
}

//GetTabbyVersion get tabbyVersion
func GetTabbyVersion() string {
	return tabbyVersion
}

//GetBuildTime get buildTime
func GetBuildTime() string {
	return BuildDate
}

//GetBuildUser get buildTime
func GetBuildUser() string {
	return BuildUser
}

//GetBuildHost get buildHost
func GetBuildHost() string {
	return BuildHostName
}

// GetHostName get host name
func GetHostName() string {
	return BuildHostName
}

//GetStartTime get start time
func GetStartTime() string {
	return startTime
}

//GetGoVersion get go version
func GetGoVersion() string {
	return goVersion
}

func PrintVersion() {
	fmt.Printf("%s : %s\n", kcolor.Green("TabbyVersion"), tabbyVersion)
	fmt.Printf("%s : %s\n", kcolor.Green("AppName"), AppName)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildTag"), BuildTag)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildVersion"), BuildVersion)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildUser"), BuildUser)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildHostname"), BuildHostName)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildStatus"), BuildStatus)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildDate"), BuildDate)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildBranch"), BuildBranch)
	fmt.Printf("%s : %s\n", kcolor.Green("BuildGitURL"), BuildGitURL)
}
