package kcrypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/khack"
	"hash/crc32"
)

// Md5Str md5 加密
// 	@Description
//	@param str
// 	@return string
func Md5Str(str string) string {
	h := md5.New()
	h.Write(khack.Slice(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Crc32
//  @Description Crc32
//  @Param str
//  @Return uint32
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE(khack.Slice(str))
}

// Sha1 sha1 加密
// 	@Description SHA1 checksum
//	@param data 要加密字符窜
// 	@return string 加密后字符内容
func Sha1(data string) string {
	s := sha1.New()
	s.Write([]byte(data))
	return hex.EncodeToString(s.Sum([]byte("")))
}
