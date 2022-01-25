// @description

package kos

import (
	"testing"
)

func TestGetHostnames(t *testing.T) {

	for _,i2 := range GetHostnames() {
		println(i2)
	}
	println("GetHostname:",GetHostname())
	println("GetOS:",GetOS())
}
