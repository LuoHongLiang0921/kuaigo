package kcrypto

import (
	"testing"
)

func TestMd5Str(t *testing.T) {
	if Md5Str("Shanghai") != "5466ee572bcbc75830d044e66ab429bc" {
		t.Fatal("should equal")
	}
}

func TestCrc32(t *testing.T) {
	if Crc32("Shanghai") != 1271261733 {
		t.Fatal("should equal")
	}
}