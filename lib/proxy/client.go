package proxy

import (
	"log"
	"net"
	"time"

	"../crypt"
	"../http"
	"../tool"
)

//Client proxy client
type Client struct {
	ProxyHost http.Host
	Listen    http.Host
	Crypter   crypt.Crypter
}

// Start start proxy client
func (client Client) Start() {

	l, err := net.Listen("tcp", client.Listen.String())
	if err != nil {
		log.Panic("端口监听失败", err)
	}

	log.Println("客户端开始监听端口:", client.Listen.String())

	for {
		brower, err := l.Accept()
		if err != nil {
			log.Panic("接受连接失败", err)
		}
		log.Println("浏览器连接成功:", brower.RemoteAddr().String())
		go client.HandleRequest(brower)
	}
}

// HandleRequest handle
func (client Client) HandleRequest(brower net.Conn) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	proxyServer, err := net.DialTimeout("tcp", client.ProxyHost.String(), time.Duration(60*time.Second))
	if err != nil {
		log.Panic("连接Proxy服务器失败: ", err)
	}
	log.Println("连接Proxy服务器成功:", client.ProxyHost.String())

	go tool.Copy(brower, proxyServer, client.Crypter)
	tool.Copy(proxyServer, brower, client.Crypter)
	brower.Close()
}
