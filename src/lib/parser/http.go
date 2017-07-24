package parser

import (
	"bytes"
	"fmt"
	"io"
	"lib/compressor"
	"lib/http"
	"lib/tool"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// HTTPParser http解析器
type HTTPParser struct {
}

// Parse http string
func (parser *HTTPParser) Parse(client io.ReadWriteCloser) (*net.TCPConn, error) {

	var buffer = make([]byte, 1024)
	length, err := client.Read(buffer)
	tool.HandleAndPanic(err, "请求数据读取流错误: ")

	var request = http.Request{}
	var method, host, protocol string
	firstLine := string(buffer[:bytes.IndexByte(buffer[:], '\n')])

	reg := regexp.MustCompile(`^(GET|POST|DELETE|PUT|CONNECT|TRACE|PATCH|HEAD)\s.*$`)

	//丢弃不规范的请求
	if len(reg.FindStringIndex(firstLine)) <= 0 {
		log.Panic("请求不规范: ", firstLine)
	}

	fmt.Sscanf(firstLine, "%s%s%s", &method, &host, &protocol)

	var re = regexp.MustCompile(`^(\w+\:\/\/)`)
	if len(re.FindStringIndex(host)) <= 0 {
		if strings.IndexAny(host, ":443") > 0 {
			host = "https://" + host
		} else {
			host = "http://" + host
		}
	}

	u, err := url.Parse(host)
	if err != nil {
		log.Println("url解析失败", err)
	}

	domain := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "80"
	}

	request.Host.Domain = domain
	request.Host.Port = port
	request.Method = method
	request.Protocol = protocol

	ip, err := net.ResolveTCPAddr("tcp", request.Host.String())
	tool.HandleAndPanic(err, "IP解析失败: ")
	webServer, err := net.DialTCP("tcp", nil, ip)
	tool.HandleAndPanic(err, "连接Web服务器失败: ", request)
	// webServer.SetKeepAlive(true)
	// webServer.SetKeepAlivePeriod(server.CheckOnline)
	// webServer.SetDeadline(time.Now().Add(server.Deadline))

	log.Println("连接Web服务器成功:", request.String())

	if request.Method == http.CONNECT {
		now := time.Now()
		// // 当请求是HTTPS请求,浏览器会发送一个CONNECT请求告诉代理服务器请求的域名和端口
		b := []byte("HTTP/1.1 200 Connection established\r\n\r\n")
		client.Write(b)
		if client, ok := client.(compressor.Compressor); ok {
			client.Flush()
		}
		log.Println("HTTPS 200 执行时间:", time.Since(now))

	} else {
		// 当请求是HTTP请求,直接传给web服务器(因为这是第一次Read读取的数据,不能漏了)
		webServer.Write(buffer[:length])
	}

	return webServer, nil
}
