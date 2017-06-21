package compressor

import "io"

// Writer 读取器
type Writer interface {
	io.Writer
	Flush() error
	Close() error
}

// Compressor 压缩器
type Compressor interface {
	NewWriter(w io.Writer, level int) (Writer, error)
	NewReader(r io.Reader) io.ReadCloser
}