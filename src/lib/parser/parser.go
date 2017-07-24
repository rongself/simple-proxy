package parser

import (
	"io"
	"net"
)

// Parser 解析器接口
type Parser interface {
	Parse(httpSteam io.ReadWriteCloser) (*net.TCPConn, error)
}
