package kcrypto

import "crypto/rc4"

// RC4
//  @Description RC4加密
//  @Param src 加密原数据
//  @Param key 加密key
//  @Return rcs加密结果
func RC4(src, key []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)
	return dst, nil
}
