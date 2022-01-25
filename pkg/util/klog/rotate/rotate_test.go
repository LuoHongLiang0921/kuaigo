// @Description

// +build linux

package rotate_test

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
