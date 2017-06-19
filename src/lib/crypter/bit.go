package crypter

// Bitcrypter 位运算加密
type Bitcrypter struct {
	Secret byte
}

// Encode 加密
func (crypter Bitcrypter) Encode(bt []byte) []byte {
	for index, value := range bt {
		bt[index] = value ^ crypter.Secret
	}
	return bt
}

// Decode 解密
func (crypter Bitcrypter) Decode(bt []byte) []byte {
	return crypter.Encode(bt)
}
