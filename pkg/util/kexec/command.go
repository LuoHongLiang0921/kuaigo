package kexec

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"os/exec"
	"runtime"
	"strings"
)

// Deprecated: 请使用 Run 替代
//  @Description 执行shell命令
//  @Param command
//  @Return data
//  @Return err
func ShellCommand(command string) (data string, err error) {
	var msg []byte
	var cmd *exec.Cmd
	cmd = exec.Command("sh", "-c", command)
	if msg, err = cmd.Output(); err != nil {
		return "", err
	}
	return string(msg), nil
}

// Command
//  @Description 执行自定义shell命令
//  @Param name
//  @Param arg
//  @Return []byte
//  @Return error
func Command(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}

// Run
//  @Description: 执行shell脚本
//  @Param arg 运行参数
//  @Param dir 指定command的工作目录，如果dir为空，则comman在调用进程所在当前目录中运行
//  @Param in 标准输入，如果stdin是nil的话，进程从null device中读取（os.DevNull），stdin也可以时一个文件，否则的话则在运行过程中再开一个goroutine去
//  @Return string
//  @Return error
func Run(arg string, dir string, in ...*bytes.Buffer) (string, error) {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case constant.OsMac, constant.OsLinux:
		cmd = exec.Command("sh", "-c", arg)
	case constant.OsWindows:
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return "", fmt.Errorf("unexpected os: %v", goos)
	}
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if len(in) > 0 {
		cmd.Stdin = in[0]
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSuffix(stderr.String(), constant.NL))
		}
		return "", err
	}

	return strings.TrimSuffix(stdout.String(), constant.NL), nil
}
