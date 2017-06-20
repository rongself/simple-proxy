package compressor

import "io"
import "compress/flate"

// FlateCompressor DEFLATE压缩器
type FlateCompressor struct {
}

// NewWriter 创建新的写入器
func (compressor *FlateCompressor) NewWriter(w io.Writer, level int) (Writer, error) {
	return flate.NewWriter(w, level)
}

// NewReader 创建新的读取器
func (compressor *FlateCompressor) NewReader(r io.Reader) io.ReadCloser {
	return flate.NewReader(r)
}
