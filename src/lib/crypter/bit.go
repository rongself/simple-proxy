package crypter

import "io"

// Bitcrypter 位运算加密
type Bitcrypter struct {
	Secret byte
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

func (crypter *Bitcrypter) Write(p []byte) (n int, err error) {
	return len(*crypter.Encode(&p)), nil
}

func (crypter *Bitcrypter) Read(p []byte) (n int, err error) {

	return len(*crypter.Encode(&p)), nil
}

func (crypter *Bitcrypter) Close() error {

	return nil
}

func (crypter *Bitcrypter) NewWriter(w io.Writer) (io.Writer, error) {

	return crypter, nil
}

func (crypter *Bitcrypter) NewReader(r io.Reader) io.ReadCloser {
	return crypter
}
