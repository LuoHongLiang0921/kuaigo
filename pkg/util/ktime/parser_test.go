// @Description

package ktime

import (
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	var testTime = Time{time.Now()}
	t1 := testTime.Parse("2020-09-05").Yesterday().Format("y-m-d")
	t2 := testTime.Parse("2020-09-05").ToFormatString("Y/m/d")
	fmt.Println("t1:",t1)
	fmt.Println("t2:",t2)
}