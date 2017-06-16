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

	l, err := net.Listen("tcp", server.Host.String())
	if err != nil {
		log.Panic("端口监听失败", err)
	}

	log.Println("开始监听端口:", server.Host.String())

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic("接受连接失败", err)
		}
		log.Println("客户端连接成功:", client.RemoteAddr().String())
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
	crypter := crypt.Bitcrypt{}
	var buffer [2048]byte
	len, err := client.Read(buffer[:])
	if err != nil {
		log.Panic("请求数据读取流错误: ", err)
	}

	request := http.Request{}
	err = request.Parse(crypter.Decode(buffer[:]))
	if err != nil {
		log.Panic("请求解析失败: ", err)
	}

	webServer, err := net.DialTimeout("tcp", request.Host.String(), time.Duration(60*time.Second))
	if err != nil {
		log.Panic("连接Web服务器失败: ", request, err)
	}

	log.Println("连接Web服务器成功:", request.String())

	if request.Method == http.CONNECT {
		client.Write(crypter.Encode([]byte("HTTP/1.1 200 Connection established\r\n\r\n")))
	} else {
		webServer.Write(buffer[:len])
	}

	go tool.Copy(webServer, client, server.Crypter)
	tool.Copy(client, webServer, server.Crypter)

	client.Close()
	webServer.Close()
}
