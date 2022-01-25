// @Description 转换格式

package kstring

import (
	"strings"
	"unicode/utf8"
)

// ToSnakeCase
//  @Description   转为snake格式
//  @Param str
//  @Return string
func ToSnakeCase(str string) string {
	str = strings.TrimSpace(strings.ToLower(str))
	return strings.Replace(str, " ", "_", -1)
}

// ToCamelCase
//  @Description  转为驼峰格式
//  @Param str
//  @Return string
func ToCamelCase(str string) string {
	str = strings.TrimSpace(str)
	if utf8.RuneCountInString(str) < 2 {
		return str
	}

	var buff strings.Builder
	var temp string
	for _, r := range str {
		c := string(r)
		if c != " " {
			if temp == " " {
				c = strings.ToUpper(c)
			}
			_, _ = buff.WriteString(c)
		}
		temp = c
	}
	return buff.String()
}
