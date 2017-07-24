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
	Host        http.Host
	Parser      parser.Parser
	Compressor  string
	Deadline    time.Duration
	CheckOnline time.Duration
}

// Start start proxy server
func (server Server) Start() {

	ip, err := net.ResolveTCPAddr("tcp", server.Host.String())
	tool.HandleAndPanic(err, "IP解析失败: ")

	proxyServer, err := net.ListenTCP("tcp", ip)
	if err != nil {
		log.Panic("服务器端口监听失败", err)
	}
	log.Println("服务器开始监听端口:", server.Host.String())
	defer proxyServer.Close()

	for {
		client, err := proxyServer.Accept()
		tool.Handle(err, "接受客户端连接失败")
		if clientTCP, ok := client.(*net.TCPConn); ok {
			log.Println("接受客户端连接成功:", client.RemoteAddr().String())
			clientTCP.SetKeepAlive(true)
			tool.Handle(err, "SetKeepAlive failed")
			clientTCP.SetKeepAlivePeriod(server.CheckOnline)
			tool.Handle(err, "SetKeepAlivePeriod failed")
			go server.HandleRequest(clientTCP)
		}
	}
}

// HandleRequest handle
func (server Server) HandleRequest(client *net.TCPConn) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer client.Close()

	var compressClientConn io.ReadWriteCloser
	if server.Compressor != "" {
		var err error
		comp, err := compressor.NewCompressor(client, flate.DefaultCompression)
		tool.HandleAndPanic(err, "初始化压缩器失败")
		compressClientConn = comp
	} else {
		compressClientConn = client
	}
	defer compressClientConn.Close()

	webServer, err := server.Parser.Parse(compressClientConn)
	tool.HandleAndPanic(err, "请求解析失败: ")
	defer webServer.Close()

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
