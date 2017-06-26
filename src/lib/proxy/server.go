package proxy

import (
	"compress/flate"
	"log"
	"net"
	"sync"
	"time"

	"io"
	"lib/compressor"
	"lib/http"
	"lib/parser"
	"lib/tool"
)

//Server proxy server
type Server struct {
	Host       http.Host
	Parser     parser.Parser
	Compressor string
	Deadline   time.Duration
}

// Start start proxy server
func (server Server) Start() {

	ip, err := net.ResolveTCPAddr("tcp", server.Host.String())
	if err != nil {
		log.Panic("IP解析失败: ", err)
	}

	proxyServer, err := net.ListenTCP("tcp", ip)
	if err != nil {
		log.Panic("服务器端口监听失败", err)
	}
	log.Println("服务器开始监听端口:", server.Host.String())
	defer proxyServer.Close()

	for {
		client, err := proxyServer.Accept()
		if err != nil {
			log.Panic("接受客户端连接失败", err)
		}

		log.Println("接受客户端连接成功:", client.RemoteAddr().String())
		go server.HandleRequest(client)
	}
}

// HandleRequest handle
func (server Server) HandleRequest(client net.Conn) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer client.Close()

	var compressClientConn io.ReadWriteCloser
	if server.Compressor != "" {
		var err error
		var compressor = &compressor.FlateCompressor{}
		compressor.Init(client, flate.DefaultCompression)
		compressClientConn = compressor
		tool.HandleAndPanic(err, "初始化压缩器失败")
	} else {
		compressClientConn = client
	}
	defer compressClientConn.Close()

	// @TOFIX 如果第一行超过了buffer长度,那就没有换行符,也就截取不到第一行,nginx 默认最大header行长度为8192byte
	var buffer = make([]byte, 2048)
	len, err := compressClientConn.Read(buffer)
	tool.HandleAndPanic(err, "请求数据读取流错误: ")

	request, err := server.Parser.Parse(buffer)
	tool.HandleAndPanic(err, "请求解析失败: ")

	webServer, err := net.DialTimeout("tcp", request.Host.String(), time.Duration(server.Deadline))
	tool.HandleAndPanic(err, "连接Web服务器失败: ", request)
	defer webServer.Close()
	webServer.SetDeadline(time.Now().Add(server.Deadline))

	log.Println("连接Web服务器成功:", request.String())

	if request.Method == http.CONNECT {
		now := time.Now()
		// // 当请求是HTTPS请求,浏览器会发送一个CONNECT请求告诉代理服务器请求的域名和端口
		b := []byte("HTTP/1.1 200 Connection established\r\n\r\n")
		compressClientConn.Write(b)
		if compressClientConn, ok := compressClientConn.(compressor.Compressor); ok {
			compressClientConn.Flush()
		}
		log.Println("HTTPS 200 执行时间:", time.Since(now))

	} else {
		// 当请求是HTTP请求,直接传给web服务器(因为这是第一次Read读取的数据,不能漏了)
		webServer.Write(buffer[:len])
	}

	var wg sync.WaitGroup

	wg.Add(1)
	//客户端过来的流量写入到目标Web服务器
	go func() {
		now := time.Now()
		w, err := tool.Copy(webServer, compressClientConn)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:proxy服务器 -> Web服务器", err)
			} else {
				log.Println("错误:proxy服务器 -> Web服务器", w, err)
			}
		}
		log.Println("proxy服务器 -> Web服务器 执行时间:", time.Since(now), "writen:", w)
		wg.Done()
	}()

	wg.Add(1)
	//目标Web服务器的相应数据写入到客户端
	go func() {
		now := time.Now()
		w, err := tool.Copy(compressClientConn, webServer)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("超时:proxy服务器 -> client", err)
			} else {
				log.Println("错误:proxy服务器 -> client", w, err)
			}

		}
		log.Println("proxy服务器 -> client 执行时间:", time.Since(now), "writen:", w)
		wg.Done()
	}()

	wg.Wait()
	log.Println("done")

}
