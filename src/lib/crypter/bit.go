package crypter

import "io"

// Bitcrypter 位运算加密
type Bitcrypter struct {
	ReadWriter io.ReadWriteCloser
	Secret     byte
}

// Encode 加密
func (crypter *Bitcrypter) Encode(bt *[]byte) *[]byte {
	for index, value := range *bt {
		(*bt)[index] = value ^ crypter.Secret
	}
	return bt
}

// Decode 解密
func (crypter *Bitcrypter) Decode(bt *[]byte) *[]byte {
	return crypter.Encode(bt)
}

//Write 写
func (crypter *Bitcrypter) Write(p []byte) (n int, err error) {
	return crypter.ReadWriter.Write(*crypter.Encode(&p))
}

//Read 读
func (crypter *Bitcrypter) Read(p []byte) (n int, err error) {
	return crypter.ReadWriter.Read(*crypter.Encode(&p))
}

//Close 关闭
func (crypter *Bitcrypter) Close() error {
	return crypter.ReadWriter.Close()
}
