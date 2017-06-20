package proxy

import (
	"compress/flate"
	"log"
	"net"

	"lib/compressor"
	"lib/crypter"
	"lib/http"
	"lib/tool"
	"time"
)

//Client proxy client
type Client struct {
	ProxyHost  http.Host
	Listen     http.Host
	Crypter    crypter.Crypter
	Compressor compressor.Compressor
}

// Start start proxy client
func (client Client) Start() {

	localServer, err := net.Listen("tcp", client.Listen.String())
	if err != nil {
		log.Panic("端口监听失败", err)
	}

	log.Println("客户端开始监听端口:", client.Listen.String())
	defer localServer.Close()

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
	defer proxyServer.Close()

	// 包装代理服务器通道
	cr := client.Compressor.NewReader(proxyServer)
	cw, err := client.Compressor.NewWriter(proxyServer, flate.DefaultCompression)
	if err != nil {
		log.Panic("初始化压缩器失败", err)
	}
	defer cr.Close()
	defer cw.Close()

	//代理过来的流量写回到浏览器
	channel := make(chan int64, 1)
	defer close(channel)
	go func() {
		w, err := tool.Copy(brower, cr, client.Crypter)
		log.Println("client -> brower", w)
		if err != nil {
			log.Println("错误:client -> brower", w, err)
		}
		channel <- w
	}()

	//浏览器过来的流量写入到代理服务器
	go func() {
		w, err := tool.Copy(cw, brower, client.Crypter)
		log.Println("brower -> proxy", w)
		if err != nil {
			log.Println("错误:brower -> proxy", w, err)
		}
		channel <- w
	}()

	w1, w2 := <-channel, <-channel
	log.Println("w1:", w1, "w2:", w2)

}
