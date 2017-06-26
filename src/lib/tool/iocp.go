package tool

import (
	"io"
	"log"

	"lib/compressor"
)

// Copy copy
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}

// CopyBuffer with buffer
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	return copyBuffer(dst, src, buf)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {

	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	for {
		nr, er := src.Read(buf)
		// log.Println(nr)
		if nr > 0 {
			b := buf[0:nr]
			nw, ew := dst.Write(b)
			if c, ok := dst.(compressor.Compressor); ok {
				err := c.Flush()
				// log.Println("Flush")
				if err != nil {
					log.Println("Flush Error:", err)
				}
			}

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
