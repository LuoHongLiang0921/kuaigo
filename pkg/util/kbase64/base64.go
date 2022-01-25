package kbase64

import "encoding/base64"

// Base64Encode
//  @Description 使用Base64编码字符串
//  @Param s 输入字符串
//  @Return 编码后的字符串
func Base64Encode(s string) string {
	return Base64EncodeByte([]byte(s))
}

// Base64EncodeByte
//  @Description 使用Base64编码字符数组
//  @Param ab 输入字符数组
//  @Return 编码后的字符串
func Base64EncodeByte(ab []byte) string {
	return base64.StdEncoding.EncodeToString(ab)
}

// Base64Decode
//  @Description 使用Base64解码字符串
//  @Param s 输入字符串
//  @Return 解码后的字符串,错误
func Base64Decode(s string) (string, error) {
	str, err := Base64DecodeByte([]byte(s))
	return string(str), err
}

// Base64DecodeByte
//  @Description 使用Base64解码字符数组
//  @Param s 输入字符数组
//  @Return 解码后的字符串,错误
func Base64DecodeByte(ab []byte) ([]byte, error) {
	str := string(ab)
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}
