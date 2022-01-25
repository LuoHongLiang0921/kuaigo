// @Description

package kmath

import (
	"fmt"
	"testing"
)

//Test
func TestFibonacci(t *testing.T) {
	for i := 0; i < 20; i++ {
		fmt.Printf("%d,", Fibonacci(i))
	}
}
