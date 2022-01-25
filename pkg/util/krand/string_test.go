// @description

package krand

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	r := GenRandomString(12, LowerUpperDigit)
	println(r)
	if len(r) != 12 {
		t.Fatal("rand str should 12 len")
	}

	t.Log(r)
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRandomString(12, LowerUpperDigit)
	}
}

func TestGenHash(t *testing.T) {
	println(GenHash("123"))
}

func TestGenRandomInt64List(t *testing.T) {
	list := GenRandomInt64List(1000,99999,100)
	for _, i2 := range list {
		fmt.Println(i2)

	}
}