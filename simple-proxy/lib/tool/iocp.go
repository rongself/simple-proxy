package tool

import (
	"io"

	"../crypt"
)

// Copy copy
func Copy(dst io.Writer, src io.Reader, crypter crypt.Crypter) (written int64, err error) {
	return copyBuffer(dst, src, nil, crypter)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte, crypter crypt.Crypter) (written int64, err error) {

	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(crypter.Encode(buf[0:nr]))
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
