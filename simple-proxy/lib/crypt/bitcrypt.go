package crypt

// Bitcrypt 位运算加密
type Bitcrypt struct {
	Secret byte
}

// Crypter 加密器
type Crypter interface {
	Encode(src []byte) []byte
	Decode(src []byte) []byte
}

// Encode 加密
func (crypter Bitcrypt) Encode(bt []byte) []byte {
	for index, value := range bt {
		bt[index] = value ^ crypter.Secret
	}
	return bt
}

// Decode 解密
func (crypter Bitcrypt) Decode(bt []byte) []byte {
	return crypter.Encode(bt)
}
