// @description

package ktime

import (
	"fmt"
	"testing"
	"time"
)

func TestTime_DaysInYear(t1 *testing.T) {
	var testTime = Time{time.Now()}
	fmt.Println("DaysInYear：",testTime.DaysInYear())

}
