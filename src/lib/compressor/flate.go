package compressor

import "io"
import "compress/flate"

// FlateCompressor DEFLATE压缩器
type FlateCompressor struct {
	reader Reader
	writer WriteFlushCloser
}

func NewCompressor(rwc io.ReadWriteCloser, level int) (*FlateCompressor, error) {
	r := flate.NewReader(rwc)
	w, err := flate.NewWriter(rwc, level)
	c := &FlateCompressor{
		reader: r,
		writer: w,
	}
	return c, err
}

//Init 初始化
func (compressor *FlateCompressor) Init(rwc io.ReadWriteCloser, level int) error {
	r := flate.NewReader(rwc)
	w, err := flate.NewWriter(rwc, level)
	compressor.reader = r
	compressor.writer = w
	return err
}

func (compressor *FlateCompressor) Write(p []byte) (n int, err error) {
	return compressor.writer.Write(p)
}

//Flush 提交数据
func (compressor *FlateCompressor) Flush() error {
	return compressor.writer.Flush()
}

func (compressor *FlateCompressor) Read(p []byte) (n int, err error) {
	return compressor.reader.Read(p)
}

// Close 关闭连接
func (compressor *FlateCompressor) Close() error {

	var err error
	if e := compressor.writer.Close(); e != nil {
		err = e
	}

	if e := compressor.writer.Close(); e != nil {
		err = e
	}

	return err
}
