// @Description 正则表达式

package kregexp

import "regexp"

// RegexpReplace
//  @Description  正则替换
//  @Param reg
//  @Param src
//  @Param temp
//  @Return string
func RegexpReplace(reg, src, temp string) string {
	var result []byte
	pattern := regexp.MustCompile(reg)
	for _, subMatches := range pattern.FindAllStringSubmatchIndex(src, -1) {
		result = pattern.ExpandString(result, temp, src, subMatches)
	}
	return string(result)
}
