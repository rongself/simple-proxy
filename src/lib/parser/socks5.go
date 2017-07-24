package parser

import (
	"io"
	"lib/compressor"
	"lib/tool"
	"log"
	"net"
	"strconv"
)

// Socks5Parser http解析器
type Socks5Parser struct {
}

// Parse http string
func (parser *Socks5Parser) Parse(client io.ReadWriteCloser) (*net.TCPConn, error) {

	var b [64]byte
	n, err := client.Read(b[:])
	tool.HandleAndPanic(err)
	if b[0] == 0x05 { //只处理Socks5协议
		//客户端回应：Socks服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})
		if comp, ok := client.(compressor.Compressor); ok {
			comp.Flush()
		}
		n, err = client.Read(b[:])
		tool.HandleAndPanic(err)
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
		ip, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
		tool.HandleAndPanic(err, "IP解析失败: ")
		webServer, err := net.DialTCP("tcp", nil, ip)
		tool.HandleAndPanic(err)
		log.Println("连接web服务器成功:", ip)
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		if comp, ok := client.(compressor.Compressor); ok {
			comp.Flush()
		}
		return webServer, nil
	}
	return nil, err
}
