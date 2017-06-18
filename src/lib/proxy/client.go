package proxy

import (
	"log"
	"net"

	"lib/crypt"

	"lib/http"
	"lib/tool"
	"time"
)

//Client proxy client
type Client struct {
	ProxyHost http.Host
	Listen    http.Host
	Crypter   crypt.Crypter
}

// Start start proxy client
func (client Client) Start() {

	localServer, err := net.Listen("tcp", client.Listen.String())
	if err != nil {
		log.Panic("端口监听失败", err)
	}

	log.Println("客户端开始监听端口:", client.Listen.String())

	for {
		brower, err := localServer.Accept()
		if err != nil {
			log.Panic("接受连接失败", err)
		}
		log.Println("浏览器连接成功:", brower.RemoteAddr().String())

		now := time.Now()
		go func() {
			client.HandleRequest(brower)
			log.Println("请求处理完成,处理时间:", time.Since(now))
		}()
	}
}

// HandleRequest handle
func (client Client) HandleRequest(brower net.Conn) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	ip, err := net.ResolveTCPAddr("tcp", client.ProxyHost.String())
	if err != nil {
		log.Panic("IP解析失败: ", err)
	}
	// 写在handle里,每一个浏览器请求建立一个TCP连接来传送
	proxyServer, err := net.DialTCP("tcp", nil, ip)
	if err != nil {
		log.Panic("连接Proxy服务器失败: ", err)
	}
	log.Println("连接Proxy服务器成功:", client.ProxyHost.String())

	go tool.Copy(brower, proxyServer, client.Crypter)
	tool.Copy(proxyServer, brower, client.Crypter)
	brower.Close()
}
