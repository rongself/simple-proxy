package compressor

import "io"

// Reader 读取器
type Reader interface {
	io.Reader
}

// Writer 写入器
type Writer interface {
	io.Writer
}

// Closer 关闭器
type Closer interface {
	io.Closer
}

// Flusher 刷新器
type Flusher interface {
	Flush() error
}

//WriteFlushCloser 组合接口
type WriteFlushCloser interface {
	Writer
	Flusher
	Closer
}

// Compressor 压缩器
type Compressor interface {
	Reader
	Writer
	Closer
	Flusher
	Init(rwc io.ReadWriteCloser, level int) error
}
