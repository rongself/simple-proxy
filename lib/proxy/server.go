package proxy

import (
	"log"
	"net"
	"time"

	"../crypt"
	"../http"
	"../tool"
)

//Server proxy server
type Server struct {
	Host    http.Host
	Crypter crypt.Crypter
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

	var buffer = make([]byte, 2048)
	len, err := client.Read(buffer[:])
	if err != nil {
		log.Panic("请求数据读取流错误: ", err)
	}

	request := http.Request{}
	err = request.Parse(server.Crypter.Decode(buffer[:]))
	if err != nil {
		log.Panic("请求解析失败: ", err)
	}

	webServer, err := net.DialTimeout("tcp", request.Host.String(), time.Duration(60*time.Second))
	if err != nil {
		log.Panic("连接Web服务器失败: ", request, err)
	}

	log.Println("连接Web服务器成功:", request.String())

	if request.Method == http.CONNECT {
		client.Write(server.Crypter.Encode([]byte("HTTP/1.1 200 Connection established\r\n\r\n")))
	} else {
		webServer.Write(buffer[:len])
	}

	go tool.Copy(webServer, client, server.Crypter)
	tool.Copy(client, webServer, server.Crypter)

	client.Close()
	webServer.Close()
}
