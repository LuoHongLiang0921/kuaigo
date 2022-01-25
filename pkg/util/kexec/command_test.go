// @Description

package kexec

import (
	"fmt"
	"os/exec"
	"testing"
)

//Test执行shell命令
func TestShellCommand(t *testing.T) {
	dateStr, _ := ShellCommand("date")
	fmt.Println(dateStr)
}

//Test执行自定义shell命令
func TestCommand(t *testing.T) {
	byte, _ := exec.Command("date", "+%Y").CombinedOutput()
	fmt.Println(string(byte))
}

//Test执行shell命令
func TestRun(t *testing.T) {
	dateStr, _ := Run("date","")

	fmt.Println(dateStr)
}
