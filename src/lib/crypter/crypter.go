package crypter

import "io"

// Crypter 加密器
type Crypter interface {
	Encode(bt *[]byte) *[]byte
	Decode(bt *[]byte) *[]byte
	NewWriter(w io.Writer) (io.Writer, error)
	NewReader(r io.Reader) io.ReadCloser
}
