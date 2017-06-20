package proxy

import (
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

	// @TOFIX 如果第一行超过了buffer长度,那就没有换行符,也就截取不到第一行,nginx 默认最大header行长度为8192byte
	var buffer = make([]byte, 2048)

	//用压缩器包装客户端
	// cr := flate.NewReader(client)
	// cw, err := flate.NewWriter(client, flate.DefaultCompression)
	// if err != nil {
	// 	log.Panic("初始化压缩器失败", err)
	// }

	len, err := client.Read(buffer)
	if err != nil {
		log.Panic("请求数据读取流错误: ", err)
	}

	// log.Println("解压后的流量:", string(buffer[:len]))

	request, err := server.Parser.Parse(buffer)
	// request, err := server.Parser.Parse(server.Crypter.Decode(buffer))
	if err != nil {
		log.Panic("请求解析失败: ", err)
	}

	webServer, err := net.DialTimeout("tcp", request.Host.String(), time.Duration(60*time.Second))
	if err != nil {
		log.Panic("连接Web服务器失败: ", request, err)
	}
	defer webServer.Close()

	log.Println("连接Web服务器成功:", request.String())

	if request.Method == http.CONNECT {
		// if err != nil {
		// 	log.Panic("初始化压缩器失败", err)
		// }
		now := time.Now()
		// // 当请求是HTTPS请求,浏览器会发送一个CONNECT请求告诉代理服务器请求的域名和端口
		client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		// // cw.Write(server.Crypter.Encode([]byte("HTTP/1.1 200 Connection established\r\n\r\n")))
		// cw.Flush()
		log.Println("HTTPS 200 执行时间:", time.Since(now))

	} else {
		// 当请求是HTTP请求,直接传给web服务器(因为这是第一次Read读取的数据,不能漏了)
		webServer.Write(buffer[:len])
	}

	channel := make(chan int64, 1)
	defer close(channel)

	//客户端过来的流量写入到目标Web服务器
	go func() {
		now := time.Now()
		w, err := tool.Copy(webServer, client, server.Crypter)
		if err != nil {
			log.Println("错误:proxy服务器 -> Web服务器", w, err)
		}
		log.Println("proxy服务器 -> Web服务器 执行时间:", time.Since(now), "writen:", w)
		channel <- w
	}()

	//目标Web服务器的相应数据写入到客户端
	go func() {
		now := time.Now()
		w, err := tool.Copy(client, webServer, server.Crypter)
		if err != nil {
			log.Println("错误:proxy服务器 -> client:", w, err)
		}
		log.Println("proxy服务器 -> client 执行时间:", time.Since(now), "writen:", w)
		channel <- w
	}()

	w1, w2 := <-channel, <-channel
	log.Println("w1", w1, "w2", w2)

}
