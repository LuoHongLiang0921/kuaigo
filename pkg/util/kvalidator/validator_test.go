// @Description

package kvalidator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsMobileNumber(t *testing.T) {

	tests := []struct {
		mobile string
		want   bool
	}{
		{"13312345678", true},
		{"14912345678", true},
		{"15312345678", true},
		{"17312345678", true},
		{"17712345678", true},
		{"18012345678", true},
		{"18112345678", true},
		{"18912345678", true},
		{"19912345678", true},
		{"13012345678", true},
		{"13112345678", true},
		{"13212345678", true},
		{"14512345678", true},
		{"15512345678", true},
		{"16612345678", true},
		{"17112345678", true},
		{"17512345678", true},
		{"17612345678", true},
		{"18512345678", true},
		{"18612345678", true},
		{"16612345678", true},
		{"13412345678", true},
		{"13512345678", true},
		{"13612345678", true},
		{"13712345678", true},
		{"13812345678", true},
		{"13912345678", true},
		{"14712345678", true},
		{"15012345678", true},
		{"15112345678", true},
		{"15212345678", true},
		{"15712345678", true},
		{"15812345678", true},
		{"15912345678", true},
		{"17212345678", true},
		{"17812345678", true},
		{"18212345678", true},
		{"18312345678", true},
		{"18412345678", true},
		{"18712345678", true},
		{"18812345678", true},
		{"19812345678", true},
	}
	for _, tt := range tests {
		t.Run(tt.mobile, func(t *testing.T) {
			if got := IsMobileNumber(tt.mobile); got != tt.want {
				t.Errorf("mobile = %s,IsMobileNumber() = %v, want %v", tt.mobile, got, tt.want)
			}
		})
	}
}

func TestURLString(t *testing.T) {
	is := assert.New(t)

	// HasURLSchema
	is.True(HasURLSchema("http://a.com"))
	is.False(HasURLSchema("abd://a.com"))
	is.False(HasURLSchema("/ab/cd"))

	// IsURL
	is.True(IsURL("a.com?p=1"))
	is.True(IsURL("http://a.com?p=1"))
	is.True(IsURL("/users/profile/1"))
	is.True(IsURL("123"))
	is.False(IsURL(""))
	
	// IsFullURL
	is.True(IsFullURL("http://a.com?p=1"))
	is.True(IsFullURL("http://www.a.com"))
	is.True(IsFullURL("https://www.a.com"))
	is.True(IsFullURL("http://a.com?p=1&c=b"))
	is.True(IsFullURL("http://a.com/ab/index"))
	is.True(IsFullURL("http://a.com/ab/index?p=1&c=b"))
	is.True(IsFullURL("http://www.a.com/ab/index?p=1&c=b"))
	is.False(IsFullURL(""))
	is.False(IsFullURL("a.com"))
	is.False(IsFullURL("a.com/ab/c"))
	is.False(IsFullURL("www.a.com"))
	is.False(IsFullURL("www.a.com?a=1"))
	is.False(IsFullURL("/users/profile/1"))
}

func TestOther(t *testing.T)  {
	is := assert.New(t)
	// IsASCII
	is.True(IsASCII("abc"))
	is.True(IsASCII("#$"))
	is.False(IsASCII(""))
	is.False(IsASCII("中文"))

	// IsMAC
	is.True(IsMAC("01:23:45:67:89:ab"))
	is.False(IsMAC("123 abc"))

	// IsRGBColor
	is.True(IsRGBColor("rgb(23,123,255)"))
	is.False(IsRGBColor(""))
	is.False(IsRGBColor("rgb(23,123,355)"))

	// IsHexColor
	is.True(IsHexColor("ccc"))
	is.True(IsHexColor("#ccc"))
	is.True(IsHexColor("ababab"))
	is.True(IsHexColor("#ababab"))
	is.False(IsHexColor(""))
}

func TestStringContains(t *testing.T) {
	// StringContains
	assert.True(t, StringContains("abc123", "123"))
	assert.False(t, StringContains("", "1234"))
	assert.False(t, StringContains("abc123", "1234"))

	// StartsWith
	assert.True(t, StartsWith("abc123", "abc"))
	assert.False(t, StartsWith("", "123"))
	assert.False(t, StartsWith("abc123", "123"))

	// EndsWith
	assert.True(t, EndsWith("abc123", "123"))
	assert.False(t, EndsWith("", "abc"))
	assert.False(t, EndsWith("abc123", "abc"))
}