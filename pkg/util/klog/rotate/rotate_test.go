// @Description
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

// +build linux

package rotate_test

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	rotate2 "git.bbobo.com/framework/tabby/pkg/util/xlog/rotate"
)

// Example of how to rotate in xresponse to SIGHUP.
func ExampleLogger_Rotate() {
	l := &rotate2.Logger{}
	log.SetOutput(l)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			l.Rotate()
		}
	}()
}
