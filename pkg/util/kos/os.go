// @Description 常用系统操作封装

package kos

import (
	"os"
	"os/user"
	"runtime"
)

// IsMacOS
//  @Description  是否是mac
//  @Return bool
func IsMacOS() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}

// IsLinux
//  @Description  是否是linux
//  @Return bool
func IsLinux() bool {
	if runtime.GOOS == "linux" {
		return true
	}
	return false
}

// IsUnix
//  @Description  是否是unix
//  @Return bool
func IsUnix() bool {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		return true
	}
	return false
}

// GetUserName
//  @Description  获取当前系统登录用户
//  @Return string
func GetUserName() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}

// GetUser
//  @Description  获取当前系统登录用户
//  @Return *user.User
func GetUser() *user.User {
	user, err := user.Current()
	if err != nil {
		return nil
	}
	return user
}

// GetHostnames
//  @Description  获取主机名多个 返回切片
//  @Return []string
func GetHostnames() []string {
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	return []string{host}
}

// GetHostname
//  @Description  获取主机名
//  @Return string
func GetHostname() string {
	hosts := GetHostnames()
	if len(hosts) == 0 {
		return "unknow"
	}
	return hosts[0]
}

// GetOS
//  @Description  获取当前操作系统
//  @Return string
func GetOS() string {
	return runtime.GOOS
}
