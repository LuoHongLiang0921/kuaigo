// @Description

package krand

import (
	"fmt"
	"testing"
	"time"
)

func TestGenUUID(t *testing.T) {
	str := GenUUID(time.Now())
	str2 := GenRandomID()
	fmt.Println(str)
	fmt.Println(str2)
}
