// @Description

package khack

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	var data []byte

	if String(data) != "" {
		t.Fatal("should be equal")
	}

	data = []byte{}
	if String(data) != "" {
		t.Fatal("should be equal")
	}

	data = []byte("elvizlai")

	if String(data) != "elvizlai" {
		t.Fatal("should be equal")
	}
}

func TestSlice(t *testing.T) {
	str := "elvizlai"

	if !reflect.DeepEqual(Slice(str), []byte(str)) {
		t.Fatal("should be equal")
	}

	str = ""

	if s := Slice(str); s != nil {
		t.Fatal("should be nil")
	}
}
