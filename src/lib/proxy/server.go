package proxy

import (
	"compress/flate"
	"log"
	"net"
	"time"

	"lib/compressor"
	"lib/crypter"
	"lib/http"
	"lib/parser"
	"lib/tool"
)

//Server proxy server
type Server struct {
	Host       http.Host
	Crypter    crypter.Crypter
	Parser     parser.Parser
	Compressor compressor.Compressor
}

// Start start proxy server
func (server Server) Start() {

	ip, err := net.ResolveTCPAddr("tcp", server.Host.String())
	if err != nil {
		log.Panic("IP解析失败: ", err)
	}

	l, err := net.ListenTCP("tcp", ip)
	if err != nil {
		log.Panic("服务器端口监听失败", err)
	}
	log.Println("服务器开始监听端口:", server.Host.String())

	for {
		client, err := l.Accept()
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

	// @TOFIX 如果第一行超过了buffer长度,那就没有换行符,也就截取不到第一行,nginx 默认最大header行长度为8192byte
	var buffer = make([]byte, 2048)

	//用压缩器包装客户端
	cr := server.Compressor.NewReader(client)
	cw, err := server.Compressor.NewWriter(client, flate.DefaultCompression)
	if err != nil {
		log.Panic("初始化压缩器失败", err)
	}

	len, err := cr.Read(buffer)
	if err != nil {
		log.Panic("请求数据读取流错误: ", err)
	}

	log.Println("解压后的流量:", string(buffer[:len]))

	request, err := server.Parser.Parse(buffer)
	// request, err := server.Parser.Parse(server.Crypter.Decode(buffer))
	if err != nil {
		log.Panic("请求解析失败: ", err)
	}

	webServer, err := net.DialTimeout("tcp", request.Host.String(), time.Duration(60*time.Second))
	if err != nil {
		log.Panic("连接Web服务器失败: ", request, err)
	}

	log.Println("连接Web服务器成功:", request.String())

	if request.Method == http.CONNECT {
		// 当请求是HTTPS请求,浏览器会发送一个CONNECT请求告诉代理服务器请求的域名和端口
		cw.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		// cw.Write(server.Crypter.Encode([]byte("HTTP/1.1 200 Connection established\r\n\r\n")))
		cw.Flush()
	} else {
		// 当请求是HTTP请求,直接传给web服务器(因为这是第一次Read读取的数据,不能漏了)
		webServer.Write(buffer[:len])
	}

	//客户端过来的流量写入到目标Web服务器
	go tool.Copy(webServer, cr, server.Crypter)

	//目标Web服务器的相应数据写入到客户端
	tool.Copy(cw, webServer, server.Crypter)
	cw.Flush()

	client.Close()
	webServer.Close()
}
