// @Description 字符串格式化

package kstring

import (
	"bytes"
	"fmt"
)

// Formatter ...
type Formatter string

// Format ...
func (fm Formatter) Format(args ...interface{}) string {
	return fmt.Sprintf(string(fm), args...)
}

// Nl2br
//  @Description  \n\r, \r\n, \r, \n 替换为 <br> 一般textarea浏览器提交带有换行符，显示的时候，需要格式化
//  @Param str
//  @Param isXhtml
//  @Return string
func Nl2br(str string, isXhtml bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if isXhtml {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}
