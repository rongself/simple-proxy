package compressor

import "io"
import "compress/zlib"

// ZlibCompressor Flate压缩器
type ZlibCompressor struct {
}

// NewWriter 创建新的写入器
func (compressor ZlibCompressor) NewWriter(w io.Writer, level int) (Writer, error) {
	return zlib.NewWriter(w), nil
}

// NewReader 创建新的读取器
func (compressor ZlibCompressor) NewReader(r io.Reader) io.ReadCloser {
	r, _ = zlib.NewReader(r)
	return nil
}
