package crypter

import "io"

// Crypter 加密器
type Crypter interface {
	Encode(bt *[]byte) *[]byte
	Decode(bt *[]byte) *[]byte
	io.ReadWriteCloser
}
