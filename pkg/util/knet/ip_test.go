// @description 

package knet

import (
	"fmt"
	"testing"
)

func TestLocalIP(t *testing.T) {
	localIp := LocalIP()
	fmt.Println("localIP:"+localIp)
}

func TestGetMacAddrs(t *testing.T) {
	macAddr := GetMacAddrs()
	fmt.Println("macAddr:",macAddr)
}

func TestGetIPs(t *testing.T) {
	ips := GetIPs()
	fmt.Println("ips:",ips)
}