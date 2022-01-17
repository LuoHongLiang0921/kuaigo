package klog

import (
	"fmt"
)

// Info ...
func Info() {
	fmt.Println("1111")
	var kzap Klog
	kzap.StartUpSugared()
}
