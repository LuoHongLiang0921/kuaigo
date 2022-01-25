// @Description 随机数据封装库

package krand

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

const Lower = 1 << 0
const Upper = 1 << 1
const Digit = 1 << 2

const LowerUpper = Lower | Upper
const LowerDigit = Lower | Digit
const UpperDigit = Upper | Digit
const LowerUpperDigit = LowerUpper | Digit

const lower = "abcdefghijklmnopqrstuvwxyz"
const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digit = "0123456789"

// RandString 随机字符串
// 	@Description
//	@param size 字符长度
//	@param set  字符类型
// 	@return string

func GenRandomString(size int, set int) string {
	charset := ""
	if set&Lower > 0 {
		charset += lower
	}
	if set&Upper > 0 {
		charset += upper
	}
	if set&Digit > 0 {
		charset += digit
	}

	lenAll := len(charset)

	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = charset[r.Intn(lenAll)]
	}
	return string(buf)
}

// GenHash
//  @Description  生成随机40位哈希
//  @Param src
//  @Return string
func GenHash(src string) string {
	// 1.获取当前时间戳
	unix := time.Now().Unix()
	// 2.将文件名和时间戳一起计算md5等到前32位十六进制字符
	hash := md5.New()
	io.WriteString(hash, src)
	io.WriteString(hash, strconv.Itoa(int(unix)))
	hb := hash.Sum(nil)

	// 获取时间戳前8位字符
	ub := strconv.Itoa(int(unix))[:8]

	// 组合输出40位哈希字符
	s := fmt.Sprintf("%x%s", hb, ub)

	return s
}
