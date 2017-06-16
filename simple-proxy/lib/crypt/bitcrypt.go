package crypt

// Secret 密钥
var Secret = []byte{0xB2, 0x09, 0xBB, 0x55, 0x93, 0x6D, 0x44, 0x47}

// Bitcrypt 位运算加密
type Bitcrypt struct {
}

// Crypter 加密器
type Crypter interface {
	Encode(src []byte) []byte
	Decode(src []byte) []byte
}

// Encode 加密
func (crypter Bitcrypt) Encode(bt []byte) []byte {
	for index, value := range bt {
		bt[index] = value ^ Secret[0]
	}
	return bt
}

// Decode 解密
func (crypter Bitcrypt) Decode(bt []byte) []byte {
	return crypter.Encode(bt)
}
