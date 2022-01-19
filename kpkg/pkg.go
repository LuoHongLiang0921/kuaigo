package kpkg

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/ktime"
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

// Name gets application name.
func Name() string {
	return AppName
}

//SetName set app anme
func SetName(s string) {
	AppName = s
}

//AppID get appID
func AppID() string {
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return BuildTag
}

//TabbyVersion get tabbyVersion
func TabbyVersion() string {
	return tabbyVersion
}

//BuildTime get buildTime
func BuildTime() string {
	return BuildDate
}

//BuildHost get buildHost
func BuildHost() string {
	return BuildHostName
}

// HostName get host name
func HostName() string {
	return BuildHostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("TabbyVersion"), kcolor.Blue(tabbyVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("AppName"), kcolor.Blue(AppName))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildTag"), kcolor.Blue(BuildTag))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildVersion"), kcolor.Blue(BuildVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildUser"), kcolor.Blue(BuildUser))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildHostname"), kcolor.Blue(BuildHostName))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildStatus"), kcolor.Blue(BuildStatus))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildDate"), kcolor.Blue(BuildDate))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildBranch"), kcolor.Blue(BuildBranch))
	fmt.Printf("%-8s]> %-30s => %s\n", "tabby", kcolor.Red("BuildGitURL"), kcolor.Blue(BuildGitURL))
}
