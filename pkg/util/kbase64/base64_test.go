package kbase64

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

// TestBase64Encode
//  @Description 测试Base64编码字符串
func TestBase64Encode(t *testing.T) {
	s := Base64Encode("abc你好")
	assert.Equal(t, s, "YWJj5L2g5aW9")
}

// TestBase64Decode
//  @Description 测试Base64解码字符串
func TestBase64Decode(t *testing.T) {
	s, _ := Base64Decode("YWJj5L2g5aW9")
	assert.Equal(t, s, "abc你好")
}
