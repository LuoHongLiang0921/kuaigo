// @Description

package kstring

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateUUID1(t *testing.T) {
	str := GenerateUUID(time.Now())
	str2 := GenerateID()
	fmt.Println(str)
	fmt.Println(str2)
}
