// @Description

package knet

import (
	"fmt"
	"net"
)

// LocalListener
//  @Description  随机一个本地端口，返回listener
//  @Return net.Listener
func LocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	return l
}
