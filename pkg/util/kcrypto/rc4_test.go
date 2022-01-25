package kcrypto

import (
	"encoding/base64"
	"testing"
)

func TestRC4(t *testing.T) {
	data, err := base64.URLEncoding.DecodeString("isEKtWvS")
	if err != nil {
		t.Fatal(err)
	}
	data, err = RC4(data, []byte("12345678"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "123abc" {
		t.Fatal("rc4 result should be 123abc")
	}
}

func BenchmarkRC4(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RC4([]byte("12345678"), []byte("123abc"))
	}
}
